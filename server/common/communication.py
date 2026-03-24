from common.utils import Bet

ACK = b'\x00'
NACK = b'\xFF'
BATCH = 0
WINNERS_REQUEST = 1

def read_exact(socket, n):
    data = bytearray()
    while len(data) < n:
        bytes_recv = socket.recv(n - len(data))
        if not bytes_recv:
            raise OSError("Connection closed by client")
        data.extend(bytes_recv)
    return data

def full_write(socket, data):
    data_sent = 0
    while data_sent < len(data):
        sent = socket.send(data[data_sent:])
        if sent == 0:
            raise OSError("Connection closed by client")
        data_sent += sent


def deserialize_into_bet(socket):
    agency = int.from_bytes(read_exact(socket, 1), byteorder='big', signed=False)

    name_len = int.from_bytes(read_exact(socket, 1), byteorder='big', signed=False)
    name = read_exact(socket, name_len).decode('utf-8')

    last_name_len = int.from_bytes(read_exact(socket, 1), byteorder='big', signed=False)
    last_name = read_exact(socket, last_name_len).decode('utf-8')

    document = int.from_bytes(read_exact(socket, 4), byteorder='big', signed=False)

    year = int.from_bytes(read_exact(socket, 2), byteorder='big', signed=False)
    month = int.from_bytes(read_exact(socket, 1), byteorder='big', signed=False)
    month = str(month) if month > 9 else "0" + str(month)
    day = int.from_bytes(read_exact(socket, 1), byteorder='big', signed=False)
    day = str(day) if day > 9 else "0" + str(day)
    birthdate = str(year) + "-" + month + "-" + day

    number = int.from_bytes(read_exact(socket, 4), byteorder='big', signed=False)

    return Bet(str(agency), name, last_name, str(document), birthdate, str(number))


def deserialize_batch(socket):
    batch_amount = int.from_bytes(read_exact(socket, 2), byteorder='big', signed=False)
    bets = []
    for _ in range(batch_amount):
        bets.append(deserialize_into_bet(socket))
    return bets

def read_client_action(socket):
    byte = read_exact(socket, 1)[0] # Get the byte from the bytearray
    action = (byte >> 7) & 0b1
    client_id = byte & 0b01111111

    if action == BATCH:
        return ('batch', client_id)
    return "winner_request", client_id

def send_ACK(socket):
    full_write(socket, ACK)

def send_NACK(socket):
    full_write(socket, NACK)


def serialize_winners(winners):
    data = bytearray()

    data.extend(len(winners).to_bytes(1, byteorder='big', signed=False))

    for winner in winners:
        data.extend(winner.to_bytes(4, byteorder='big', signed=False))

    return data

def send_winners(socket, winners):
    payload = serialize_winners(winners)
    full_write(socket, payload)

