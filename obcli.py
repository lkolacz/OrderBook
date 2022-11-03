import asyncio
import sys
import websockets
from enum import Enum

HOST = "127.0.0.1"
PORT = "8007"


class TerminalStyle(str, Enum):
    ENDA = '\033[0;0m'
    ENDC = '\033[0m'
    GREEN = '\033[92m'
    RED = '\033[41m'
    START = '\033[1m'


def get_websocket_endpoint(host=HOST, port=PORT):
    return f'ws://{host or HOST}:{port or PORT}/api/v1/order-handling'


async def handle_action():
    ep_addr = get_websocket_endpoint()
    async with websockets.connect(ep_addr) as websocket:
        print(f"{TerminalStyle.GREEN}Enter JSON transaction (or prompt empty to exit):{TerminalStyle.ENDC}")
        while True:
            try:
                command = input()
                if command != '':
                    await websocket.send(command)
                    try:
                        while True:
                            receive = await asyncio.wait_for(websocket.recv(), timeout=1.0)
                            sys.stdout.write(f"{receive}")
                    except asyncio.exceptions.TimeoutError:
                        pass

                else:
                    print(f"{TerminalStyle.RED}The OrderBook Client is going down... (in progress){TerminalStyle.ENDC}")
                    sys.exit()
            except Exception as err:
                sys.stderr.write("Disconnected from OrderBook server due to error: {err}. Reconnecting... ")
                websocket = await websockets.connect(ep_addr)

async_loop = asyncio.get_event_loop()


while True:
    try:
        async_loop.run_until_complete(handle_action())
    except ConnectionRefusedError:
        print('Could not connect to server. Is it running?')
        break
