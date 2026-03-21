import socket
import logging
import signal
from common.utils import store_bets
from common.communication import deserialize_into_bet, full_write, send_ACK, send_NACK


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
            bet = deserialize_into_bet(client_sock)
            store_bets([bet])
            logging.info(f'action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}')
            send_ACK(client_sock)
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            send_NACK(client_sock)
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

    def __graceful_exit(self, signum, _):
        print("Gracefully shuting down server")
        print(f"SIGNAL: {signum}")
        self._exit = True


