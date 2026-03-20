import socket
import logging
import signal
from common.utils import Bet, store_bets


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_socket.settimeout(1)
        self._exit = False
        signal.signal(signal.SIGTERM, self.__graceful_exit)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while not self._exit:
            client_sock = self.__accept_new_connection()
            if client_sock == None:
                continue
            self.__handle_client_connection(client_sock)

        self._server_socket.close()

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # TODO: Modify the receive to avoid short-reads
            bet = self.__read_msg(client_sock)
            store_bets([bet])
            logging.info(f'action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}')
            # TODO: Modify the send to avoid short-writes
            client_sock.send("{}\n".format("Recibi el mensaje").encode('utf-8'))
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        logging.info('action: accept_connections | result: in_progress')
        while not self._exit:
            try:
                c, addr = self._server_socket.accept()
                logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
                # Connection arrived
                return c
            except socket.timeout:
                continue

        # SIGTERM receive
        return None

    def __read_msg(self, client_sock):
        agency = int.from_bytes(client_sock.recv(1), byteorder='big', signed=False)

        name_len = int.from_bytes(client_sock.recv(4), byteorder='big', signed=False)
        name = client_sock.recv(name_len).rstrip().decode('utf-8')

        last_name_len = int.from_bytes(client_sock.recv(4), byteorder='big', signed=False)
        last_name = client_sock.recv(last_name_len).rstrip().decode('utf-8')

        document = int.from_bytes(client_sock.recv(4), byteorder='big', signed=False)

        year = int.from_bytes(client_sock.recv(1), byteorder='big', signed=False)
        month = int.from_bytes(client_sock.recv(1), byteorder='big', signed=False)
        day = int.from_bytes(client_sock.recv(1), byteorder='big', signed=False)
        birthdate = str(year) + "-" + str(month) + "-" + str(day)

        number = int.from_bytes(client_sock.recv(4), byteorder='big', signed=False)

        print(str(agency), name, last_name, str(document), birthdate, str(number))

        return Bet(str(agency), name, last_name, str(document), birthdate, str(number))


    def __graceful_exit(self, signum, _):
        print("Gracefully shuting down server")
        print(f"SIGNAL: {signum}")
        self._exit = True
