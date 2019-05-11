import socket
import time


def send_msg(sock, msg):
    try:
        sock[0].sendall(msg)
        data = sock[0].recv(1024)
        print(sock[1]+":", repr(data), len(data))
    except socket.timeout:
        print(f"{sock[1]} timed out")


def recv_msg(sock):
    try:
        data = sock[0].recv(1024)
        print(sock[1]+":", repr(data), len(data))
    except socket.timeout:
        print(f"{sock[1]} timed out waiting for a message")


if __name__ == '__main__':
    host = 'localhost'
    port = 28813

    sock1 = (socket.socket(socket.AF_INET, socket.SOCK_STREAM), "client1")
    sock1[0].connect((host, port))
    sock1[0].settimeout(1)

    sock2 = (socket.socket(socket.AF_INET, socket.SOCK_STREAM), "client2")
    sock2[0].connect((host, port))
    sock2[0].settimeout(1)

    send_msg(sock1, b'saddfjalksjd f\r\n')
    send_msg(sock1, b'dslp/2.0\r\nrequest time\r\ndslp/body\r\n')
    send_msg(sock2, b'dslp/2.0\r\nrequest time\r\ndslp/body\r\n')
    send_msg(sock1, b'dslp/2.0\r\nuser join\r\nclient1\r\ndslp/body\r\n')
    send_msg(sock2, b'dslp/2.0\r\nuser join\r\nclient2\r\ndslp/body\r\n')
    send_msg(sock2, b'dslp/2.0\r\nuser text notify\r\nclient2\r\nclient1\r\n1\r\ndslp/body\r\nHello client1!\r\n')
    recv_msg(sock1)
    send_msg(sock1, b'dslp/2.0\r\nuser file notify\r\nclient1\r\nclient2\r\ntest.txt\r\nidk/text\r\n8\r\ndslp/body\r\n12345678')
    recv_msg(sock2)
