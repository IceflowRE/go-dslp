import socket
import time

if __name__ == '__main__':
    host = 'localhost'
    port = 28813

    total_data = []
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        sock.connect((host, port))
        #sock.sendall(b'saddfjalksjd f\r\n')
        #sock.sendall(b'dslp/2.0\r\nrequest time\r\ndslp/body\r\n')
        #data = sock.recv(1024)
        #print(repr(data), len(data))
        #time.sleep(1)
        #sock.sendall(b'dslp/2.0\r\njoin group\r\nhello world\r\ndslp/body\r\n')
        time.sleep(1)
        sock.sendall(b'dslp/2.0\r\ngroup notify\r\nhello world\r\n1\r\ndslp/body\r\nhello world!\r\n')
        data = sock.recv(1024)
        print(repr(data), len(data))
