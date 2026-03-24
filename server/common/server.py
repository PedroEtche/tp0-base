import socket
import logging
import signal
from common.utils import store_bets, load_bets, has_won
from common.communication import deserialize_batch, read_client_action, send_ACK, send_NACK, send_winners


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_socket.settimeout(1)
        # Clients config. Use to know how many clients will notify bets and when to start the giveaway process
        self._clients_amount = listen_backlog
        self._clients_listen = 0
        # Set SIGTERM resolution
        signal.signal(signal.SIGTERM, self.__graceful_exit)
        self._exit = False

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
        client_action, client_id = read_client_action(client_sock)
        if client_action == 'batch':
            self.__handle_batch(client_sock)
        else:
            self.__handle_giveaway_request(client_sock, client_id)

        client_sock.close()



    def __handle_batch(self, client_sock):
        bets = []
        while True: 
            try:
                bets = deserialize_batch(client_sock)
                store_bets(bets)
                logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
                send_ACK(client_sock)
                bets = []
            except Exception as e:
                if isinstance(e, OSError) and str(e) == "Connection closed by client":
                    logging.info('action: connection_closed | result: success')
                    break
                else: 
                    logging.error(f'action: apuesta_recibida | result: fail | cantidad: {len(bets)}')
                    send_NACK(client_sock)

        self._clients_listen += 1
        if self._clients_listen == self._clients_amount:
            logging.info('action: sorteo | result: success')



    def __handle_giveaway_request(self, client_sock, client_id):
        try:
            # Check if all clients have notified their bets. If not, send NACK and return
            if self._clients_listen != self._clients_amount:
                send_NACK(client_sock)
                return

            bets = load_bets()
            winners = []
            for bet in bets:
                if bet.agency == str(client_id) and has_won(bet):
                    winners.append(int(bet.document))

            send_winners(client_sock, winners)
        except Exception as e:
            logging.error(f'action: pedido_ganadores | result: fail | err: {e}')

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

    def __graceful_exit(self, signum, _):
        print("Gracefully shuting down server")
        print(f"SIGNAL: {signum}")
        self._exit = True


