import socket
import time

if __name__ == '__main__':
    host = 'localhost'
    port = 28813

    total_data = []
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        sock.connect((host, port))
        sock.sendall(b'saddfjalksjd f\r\n')
        sock.sendall(b'dslp/1.2\r\nrequest time\r\ndslp/end\r\n')
        data = sock.recv(1024)
        print(repr(data), len(data))
        time.sleep(1)
        sock.sendall(b'dslp/1.2\r\nrequest time\r\ndslp/end\r\n')
        data = sock.recv(1024)
        print(repr(data), len(data))
        time.sleep(1)
        sock.sendall(b'dslp/1.2\r\ngroup join\r\ntest\r\ndslp/end\r\n')
        time.sleep(1)
        sock.sendall(b'dslp/1.2\r\ngroup notify\r\ntest\r\nHello World!\r\ndslp/end\r\n')
        data = sock.recv(1024)
        print(repr(data), len(data))
        sock.sendall(b'dslp/1.2\r\nrequest time\r\ndslp/end\r\n')
        data = sock.recv(1024)
        print(repr(data), len(data))
        sock.sendall(b'dslp/1.2\r\npeer notify\r\n127.0.0.1\r\n127.0.0.1\r\ndslp/end\r\n')
        data = sock.recv(1024)
        print(repr(data), len(data))
        #sock.sendall(b'dslp/1.2\r\npeer notfiy\r\n127.0.0.1\r\ndslp/end\r\n')
        #data = sock.recv(1024)
        #print(repr(data), len(data))
