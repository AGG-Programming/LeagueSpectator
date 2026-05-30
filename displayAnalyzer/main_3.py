import cv2
import mss
import pytesseract
import re
import websockets
import asyncio
import threading
import queue
import json
import time
import numpy as np

# Tesseract Pfad-Konfiguration
pytesseract.pytesseract.tesseract_cmd = r'C:\Program Files\Tesseract-OCR\tesseract.exe'

# --- DEINE BBOX DEFINITIONEN ---
BBOX_BLUE_TEAM_GOLD = (800, 13, 851, 39)
BBOX_RED_TEAM_GOLD = (1105, 14, 1156, 39)

BBOX_BLUE_GOLD_LIST = [(650, 865, 800, 890), (650, 910, 800, 935), (650, 950, 800, 975), (650, 995, 800, 1020), (650, 1040, 800, 1065)]
BBOX_RED_GOLD_LIST = [(1130, 865, 1290, 890), (1130, 910, 1290, 935), (1130, 950, 1290, 975), (1130, 995, 1290, 1020), (1130, 1040, 1290, 1065)]

BBOX_BLUE_BOUNTY_LIST = [(820, 847, 860, 870), (820, 893, 860, 915), (820, 937, 860, 957), (820, 981, 860, 1001), (820, 1025, 860, 1045)]
BBOX_RED_BOUNTY_LIST = [(1076, 847, 1118, 870), (1076, 893, 1118, 915), (1076, 937, 1118, 957), (1076, 981, 1118, 1001), (1076, 1025, 1118, 1045)]

BBOX_BLUE_CREEP_LIST = [(877, 860, 915, 885), (877, 905, 915, 930), (877, 950, 915, 975), (877, 995, 915, 1020), (877, 1040, 915, 1065)]
BBOX_RED_CREEP_LIST = [(1018, 860, 1055, 885), (1018, 905, 1055, 930), (1018, 950, 1055, 975), (1018, 995, 1055, 1020), (1018, 1040, 1055, 1065)]

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
            return cv2.cvtColor(np.array(screenshot), cv2.COLOR_BGRA2BGR)
        return None

def createMask(image, mask_type):
    hsv = cv2.cvtColor(image, cv2.COLOR_BGR2HSV)
    if mask_type == "blue":
        return cv2.inRange(hsv, np.array([90, 100, 100]), np.array([120, 255, 255]))
    elif mask_type == "red":
        mask1 = cv2.inRange(hsv, np.array([0, 150, 80]), np.array([10, 255, 255]))
        mask2 = cv2.inRange(hsv, np.array([170, 150, 80]), np.array([180, 255, 255]))
        return cv2.bitwise_or(mask1, mask2)
    elif mask_type == "white":
        return cv2.inRange(hsv, np.array([0, 0, 180]), np.array([180, 50, 255]))
    elif mask_type == "gold":
        return cv2.inRange(hsv, np.array([15, 60, 110]), np.array([30, 200, 255]))
    return None

def preprocessImage(image, mask_type):
    if image is None: return None
    resized = cv2.resize(image, None, fx=1.5, fy=1.5, interpolation=cv2.INTER_CUBIC)
    mask = createMask(resized, mask_type)
    if mask is None: return None
    return cv2.GaussianBlur(mask, (3, 3), 0)

def parseGold(gold_value):
    matches = re.findall(r'\d+', gold_value)
    return [int(num) for num in matches]

def processCategoryBatch(bbox_list, mask_type):
    """
    Schneidet alle 5 Boxen aus, verarbeitet sie vor, klebt sie vertikal 
    zusammen und jagt sie als EINEN Block durch Tesseract (Massiver Boost!)
    """
    processed_images = []
    
    for bbox in bbox_list:
        cropped = captureScreenArea(bbox, monitor_index=3)
        processed = preprocessImage(cropped, mask_type)
        if processed is not None:
            # Füge einen kleinen schwarzen Rand zwischen den Bildern hinzu, 
            # damit Tesseract die Zeilen sauber trennen kann
            bordered = cv2.copyMakeBorder(processed, 4, 4, 0, 0, cv2.BORDER_CONSTANT, value=0)
            processed_images.append(bordered)
            
    if not processed_images:
        return ["0"] * 5

    # Die 5 Bilder vertikal zu einer langen Spalte zusammenkleben
    vertical_strip = cv2.vconcat(processed_images)
    
    # PSM 6 (Einheitlicher Textblock) ist hier perfekt, da wir ein mehrzeiliges Dokument simulieren
    customConfig = r'--oem 3 --psm 6 -c tessedit_char_whitelist=0123456789.kK()G'
    raw_text = pytesseract.image_to_string(vertical_strip, config=customConfig)
    
    # Text zeilenweise aufsplitten und leere Einträge filtern
    lines = [line.strip() for line in raw_text.split('\n') if line.strip()]
    
    # Sichern, dass wir immer genau 5 Elemente zurückgeben (Fallback bei Erkennungsfehlern)
    while len(lines) < 5:
        lines.append("0")
    return lines[:5]

def getData():
    """Sammelt alle Spieldaten über optimiertes Batch-OCR"""
    
    # 1. Team-Gold auslesen (Einzelboxen, aber schnelles PSM 13)
    customConfig_single = r'--oem 3 --psm 13 -c tessedit_char_whitelist=0123456789.kK()'
    
    blue_team_raw = pytesseract.image_to_string(preprocessImage(captureScreenArea(BBOX_BLUE_TEAM_GOLD, 3), "blue"), config=customConfig_single)
    red_team_raw = pytesseract.image_to_string(preprocessImage(captureScreenArea(BBOX_RED_TEAM_GOLD, 3), "red"), config=customConfig_single)
    
    # 2. Spieler-Batches verarbeiten (Jeder Aufruf liest 5 Spieler auf einmal!)
    blue_gold_lines = processCategoryBatch(BBOX_BLUE_GOLD_LIST, "white")
    red_gold_lines = processCategoryBatch(BBOX_RED_GOLD_LIST, "white")
    
    blue_creep_lines = processCategoryBatch(BBOX_BLUE_CREEP_LIST, "gold")
    red_creep_lines = processCategoryBatch(BBOX_RED_CREEP_LIST, "gold")
    
    blue_bounty_lines = processCategoryBatch(BBOX_BLUE_BOUNTY_LIST, "gold")
    red_bounty_lines = processCategoryBatch(BBOX_RED_BOUNTY_LIST, "gold")
    
    # 3. JSON-Struktur aufbauen
    data = {
        "gold": {"blue": {}, "red": {}},
        "creep": {"blue": {}, "red": {}},
        "bounty": {"blue": {}, "red": {}},
        "team_gold": {"blue_team": 0, "red_team": 0}
    }
    
    # Daten einsortieren und KDA-Gold parsen
    for i in range(5):
        key_b = f"b{i+1}"
        key_r = f"r{i+1}"
        
        bg_vals = parseGold(blue_gold_lines[i])
        rg_vals = parseGold(red_gold_lines[i])
        
        data["gold"]["blue"][key_b] = bg_vals[1] if len(bg_vals) >= 2 else 0
        data["gold"]["red"][key_r] = rg_vals[1] if len(rg_vals) >= 2 else 0
        
        data["creep"]["blue"][key_b] = blue_creep_lines[i]
        data["creep"]["red"][key_r] = red_creep_lines[i]
        
        data["bounty"]["blue"][key_b] = blue_bounty_lines[i]
        data["bounty"]["red"][key_r] = red_bounty_lines[i]
        
    # Teamgold aus den Einzelwerten summieren
    data["team_gold"]["blue_team"] = sum(data["gold"]["blue"].values())
    data["team_gold"]["red_team"] = sum(data["gold"]["red"].values())
    
    return data

# --- WEBSOCKET SERVER LOGIK ---
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
        
        # Daten abrufen
        data = getData()
        
        print(f"Analyse-Dauer: {time.time() - start_time:.2f} Sekunden")
        
        if CONNECTED_CLIENTS:
            DATA_QUEUE.put(data)