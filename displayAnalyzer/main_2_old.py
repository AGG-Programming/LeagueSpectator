import cv2
import mss
import pytesseract
import re
import websockets
import asyncio
import threading
import queue
import pydirectinput
import json
import time
import numpy as np
from concurrent.futures import ThreadPoolExecutor  # NEU: Für die Parallelisierung

# Tesseract Pfad-Konfiguration
pytesseract.pytesseract.tesseract_cmd = r'C:\Program Files\Tesseract-OCR\tesseract.exe'

# Globaler ThreadPool mit maximal 5 Workern
EXECUTOR = ThreadPoolExecutor(max_workers=5)

# --- DEINE BBOX DEFINITIONEN (unverändert) ---
BBOX_BLUE_TEAM_GOLD = (800, 13, 851, 39)
BBOX_RED_TEAM_GOLD = (1105, 14, 1156, 39)

BBOX_B1_GOLD = (650, 865, 800, 890); BBOX_B2_GOLD = (650, 910, 800, 935); BBOX_B3_GOLD = (650, 950, 800, 975); BBOX_B4_GOLD = (650, 995, 800, 1020); BBOX_B5_GOLD = (650, 1040, 800, 1065)
BBOX_R1_GOLD = (1130, 865, 1290, 890); BBOX_R2_GOLD = (1130, 910, 1290, 935); BBOX_R3_GOLD = (1130, 950, 1290, 975); BBOX_R4_GOLD = (1130, 995, 1290, 1020); BBOX_R5_GOLD = (1130, 1040, 1290, 1065)
BBOX_B1_BOUNTY = (820, 847, 860, 870); BBOX_B2_BOUNTY = (820, 893, 860, 915); BBOX_B3_BOUNTY = (820, 937, 860, 957); BBOX_B4_BOUNTY = (820, 981, 860, 1001); BBOX_B5_BOUNTY = (820, 1025, 860, 1045)
BBOX_R1_BOUNTY = (1076, 847, 1118, 870); BBOX_R2_BOUNTY = (1076, 893, 1118, 915); BBOX_R3_BOUNTY = (1076, 937, 1118, 957); BBOX_R4_BOUNTY = (1076, 937, 1118, 957); BBOX_R5_BOUNTY = (1076, 1025, 1118, 1045)
BBOX_B1_CREEP = (877, 860, 915, 885); BBOX_B2_CREEP = (877, 905, 915, 930); BBOX_B3_CREEP = (877, 950, 915, 975); BBOX_B4_CREEP = (877, 995, 915, 1020); BBOX_B5_CREEP = (877, 1040, 915, 1065)
BBOX_R1_CREEP = (1018, 860, 1055, 885); BBOX_R2_CREEP = (1018, 905, 1055, 930); BBOX_R3_CREEP = (1018, 950, 1055, 975); BBOX_R4_CREEP = (1018, 995, 1055, 1020); BBOX_R5_CREEP = (1018, 1040, 1055, 1065)

# --- STRUKTURIERTE LISTEN FÜR DIE PARALLELE VERARBEITUNG ---
# Jedes Element enthält zusätzlich einen eindeutigen Identifikations-Key für das finale JSON
tasks_config = [
    # Blue Player Gold
    {"bbox": BBOX_B1_GOLD, "mask_type": "white", "category": "gold_blue", "key": "b1"},
    {"bbox": BBOX_B2_GOLD, "mask_type": "white", "category": "gold_blue", "key": "b2"},
    {"bbox": BBOX_B3_GOLD, "mask_type": "white", "category": "gold_blue", "key": "b3"},
    {"bbox": BBOX_B4_GOLD, "mask_type": "white", "category": "gold_blue", "key": "b4"},
    {"bbox": BBOX_B5_GOLD, "mask_type": "white", "category": "gold_blue", "key": "b5"},
    
    # Red Player Gold
    {"bbox": BBOX_R1_GOLD, "mask_type": "white", "category": "gold_red", "key": "r1"},
    {"bbox": BBOX_R2_GOLD, "mask_type": "white", "category": "gold_red", "key": "r2"},
    {"bbox": BBOX_R3_GOLD, "mask_type": "white", "category": "gold_red", "key": "r3"},
    {"bbox": BBOX_R4_GOLD, "mask_type": "white", "category": "gold_red", "key": "r4"},
    {"bbox": BBOX_R5_GOLD, "mask_type": "white", "category": "gold_red", "key": "r5"},
    
    # Blue Player Creep
    {"bbox": BBOX_B1_CREEP, "mask_type": "gold", "category": "creep_blue", "key": "b1"},
    {"bbox": BBOX_B2_CREEP, "mask_type": "gold", "category": "creep_blue", "key": "b2"},
    {"bbox": BBOX_B3_CREEP, "mask_type": "gold", "category": "creep_blue", "key": "b3"},
    {"bbox": BBOX_B4_CREEP, "mask_type": "gold", "category": "creep_blue", "key": "b4"},
    {"bbox": BBOX_B5_CREEP, "mask_type": "gold", "category": "creep_blue", "key": "b5"},
    
    # Red Player Creep
    {"bbox": BBOX_R1_CREEP, "mask_type": "gold", "category": "creep_red", "key": "r1"},
    {"bbox": BBOX_R2_CREEP, "mask_type": "gold", "category": "creep_red", "key": "r2"},
    {"bbox": BBOX_R3_CREEP, "mask_type": "gold", "category": "creep_red", "key": "r3"},
    {"bbox": BBOX_R4_CREEP, "mask_type": "gold", "category": "creep_red", "key": "r4"},
    {"bbox": BBOX_R5_CREEP, "mask_type": "gold", "category": "creep_red", "key": "r5"},
    
    # Blue Player Bounty
    {"bbox": BBOX_B1_BOUNTY, "mask_type": "gold", "category": "bounty_blue", "key": "b1"},
    {"bbox": BBOX_B2_BOUNTY, "mask_type": "gold", "category": "bounty_blue", "key": "b2"},
    {"bbox": BBOX_B3_BOUNTY, "mask_type": "gold", "category": "bounty_blue", "key": "b3"},
    {"bbox": BBOX_B4_BOUNTY, "mask_type": "gold", "category": "bounty_blue", "key": "b4"},
    {"bbox": BBOX_B5_BOUNTY, "mask_type": "gold", "category": "bounty_blue", "key": "b5"},
    
    # Red Player Bounty
    {"bbox": BBOX_R1_BOUNTY, "mask_type": "gold", "category": "bounty_red", "key": "r1"},
    {"bbox": BBOX_R2_BOUNTY, "mask_type": "gold", "category": "bounty_red", "key": "r2"},
    {"bbox": BBOX_R3_BOUNTY, "mask_type": "gold", "category": "bounty_red", "key": "r3"},
    {"bbox": BBOX_R4_BOUNTY, "mask_type": "gold", "category": "bounty_red", "key": "r4"},
    {"bbox": BBOX_R5_BOUNTY, "mask_type": "gold", "category": "bounty_red", "key": "r5"},
]

def captureScreenArea(bbox, monitor_index=3):
    with mss.mss() as sct:
        if monitor_index < len(sct.monitors):
            monitor = sct.monitors[monitor_index]
            x1, y1, x2, y2 = bbox
            cfg = {
                "top": monitor["top"] + y1,
                "left": monitor["left"] + x1,
                "width": x2 - x1,
                "height": y2 - y1
            }
            screenshot = sct.grab(cfg)
            frame = np.array(screenshot)
            return cv2.cvtColor(frame, cv2.COLOR_BGRA2BGR)
        return None

def createMask(image, mask_type):
    hsv = cv2.cvtColor(image, cv2.COLOR_BGR2HSV)
    if mask_type == "blue":
        lower_blue = np.array([90, 100, 100])
        upper_blue = np.array([120, 255, 255])
        return cv2.inRange(hsv, lower_blue, upper_blue)
    elif mask_type == "red":
        lower_red1 = np.array([0, 150, 80])
        upper_red1 = np.array([10, 255, 255])
        lower_red2 = np.array([170, 150, 80])
        upper_red2 = np.array([180, 255, 255])
        mask1 = cv2.inRange(hsv, lower_red1, upper_red1)
        mask2 = cv2.inRange(hsv, lower_red2, upper_red2)
        return cv2.bitwise_or(mask1, mask2)
    elif mask_type == "white":
        lower_white = np.array([0, 0, 180])
        upper_white = np.array([180, 50, 255])
        return cv2.inRange(hsv, lower_white, upper_white)
    elif mask_type == "gold":
        lower_gold = np.array([15, 60, 110])
        upper_gold = np.array([30, 200, 255])
        return cv2.inRange(hsv, lower_gold, upper_gold)
    return None

def preprocessImage(image, mask_type):
    if image is None:
        return None
    resized = cv2.resize(image, None, fx=1.5, fy=1.5, interpolation=cv2.INTER_CUBIC)
    mask = createMask(resized, mask_type)
    if mask is None: return None
    processed = cv2.GaussianBlur(mask, (3, 3), 0)
    return processed

def getNumbersFromScreen(bbox, mask_type):
    cropped = captureScreenArea(bbox, monitor_index=3)
    if cropped is None:
        return "Error (No Monitor)"
        
    processed = preprocessImage(cropped, mask_type)
    
    # Optisches Debug-Fenster (Achtung: cv2.imshow in Multithreading kann flackern)
    # cv2.imshow("Debug", processed)
    # cv2.waitKey(1)
    
    # PSM 13 ist deutlich schneller für einzelne, kurze Zeilenkonstrukte!
    customConfig = r'--oem 3 --psm 13 -c tessedit_char_whitelist=0123456789.kK()G'
    text = pytesseract.image_to_string(processed, config=customConfig)
    return text.strip()

def parseGold(gold_value):
    matches = re.findall(r'\d+', gold_value)
    return [int(num) for num in matches]

def workerTask(task):
    """Verarbeitungseinheit für ein einzelnes Bbox-Element im ThreadPool"""
    raw_value = getNumbersFromScreen(task["bbox"], task["mask_type"])
    return {
        "category": task["category"],
        "key": task["key"],
        "raw_value": raw_value
    }

def calcTotalTeamGold(blue_team_gold, red_team_gold):
    blue_team_total = sum(blue_team_gold.values())
    red_team_total = sum(red_team_gold.values())
    return {"blue_team": blue_team_total, "red_team": red_team_total}

def getData():
    """Nutzt den ThreadPoolExecutor, um alle Boxen parallel zu scannen"""
    # Alle 30 Tasks gleichzeitig an den Thread-Pool übergeben
    results = list(EXECUTOR.map(workerTask, tasks_config))
    
    # Zielstruktur aufbauen
    data = {
        "gold": {"blue": {}, "red": {}},
        "creep": {"blue": {}, "red": {}},
        "bounty": {"blue": {}, "red": {}},
        "team_gold": {"blue_team": 0, "red_team": 0}
    }
    
    # Ergebnisse sortiert einsortieren
    for res in results:
        cat = res["category"]
        key = res["key"]
        val = res["raw_value"]
        
        if cat == "gold_blue":
            gold_values = parseGold(val)
            # Sichert ab, dass das KDA-Feld mindestens 2 Zahlenwerte hatte (z.B. Kills/Deaths)
            data["gold"]["blue"][key] = gold_values[1] if len(gold_values) >= 2 else 0
        elif cat == "gold_red":
            gold_values = parseGold(val)
            data["gold"]["red"][key] = gold_values[1] if len(gold_values) >= 2 else 0
        elif cat == "creep_blue":
            data["creep"]["blue"][key] = val
        elif cat == "creep_red":
            data["creep"]["red"][key] = val
        elif cat == "bounty_blue":
            data["bounty"]["blue"][key] = val
        elif cat == "bounty_red":
            data["bounty"]["red"][key] = val

    # Teamgold berechnen
    data["team_gold"] = calcTotalTeamGold(data["gold"]["blue"], data["gold"]["red"])
    return data

# --- WEBSOCKET SERVER LOGIK (unverändert) ---
DATA_QUEUE = queue.Queue()
CONNECTED_CLIENTS = set()

async def register(websocket):
    CONNECTED_CLIENTS.add(websocket)
    print(f"New client connected. Active connections: {len(CONNECTED_CLIENTS)}")
    try:
        await websocket.wait_closed()
    finally:
        CONNECTED_CLIENTS.remove(websocket)
        print(f"Client disconnected. Active connections: {len(CONNECTED_CLIENTS)}")

async def broadcastFromQueue():
    while True:
        loop = asyncio.get_running_loop()
        try:
            data_dict = await loop.run_in_executor(None, lambda: DATA_QUEUE.get(timeout=0.5))
            if CONNECTED_CLIENTS and data_dict:
                json_message = json.dumps(data_dict)
                current_clients = list(CONNECTED_CLIENTS)
                for client in current_clients:
                    try:
                        await client.send(json_message)
                    except Exception as send_error:
                        print(f"Error while sending to client: {send_error}")
                print(f"Data successfully sent to {len(current_clients)} client(s).")
        except queue.Empty:
            pass
        await asyncio.sleep(0.01)

def startWebsocketServer():
    async def main():
        async with websockets.serve(register, "localhost", 8765):
            print("WebSocket-Server runs on ws://localhost:8765")
            await broadcastFromQueue()
    asyncio.run(main())

if __name__ == "__main__":
    print("Script active. Scanning MONITOR 2 (Index 3 in MSS) for League of Legends...")
    
    server_thread = threading.Thread(target=startWebsocketServer, daemon=True)
    server_thread.start()

    time.sleep(1)

    while True:
        start_time = time.time()
        
        # Holt alle Daten parallel über den ThreadPool
        data = getData()
        
        print(f"Analyse-Dauer: {time.time() - start_time:.2f} Sekunden")
        print(f"data: {data}")
        
        if not CONNECTED_CLIENTS:
            print("Data gathered but no client connected")
        else:
            DATA_QUEUE.put(data)