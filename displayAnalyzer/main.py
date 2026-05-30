import cv2
import mss
import pytesseract
import re
import pydirectinput
import time
import numpy as np

# Tesseract Pfad-Konfiguration
pytesseract.pytesseract.tesseract_cmd = r'C:\Program Files\Tesseract-OCR\tesseract.exe'

BBOX_BLUE_TEAM_GOLD = (800, 13, 851, 39)
BBOX_RED_TEAM_GOLD = (1105, 14, 1156, 39)

BBOX_B1_GOLD = (650, 865, 800, 890)
BBOX_B2_GOLD = (650, 910, 800, 935)
BBOX_B3_GOLD = (650, 950, 800, 975)
BBOX_B4_GOLD = (650, 995, 800, 1020)
BBOX_B5_GOLD = (650, 1040, 800, 1065)
BBOX_R1_GOLD = (1130, 865, 1290, 890)
BBOX_R2_GOLD = (1130, 910, 1290, 935)
BBOX_R3_GOLD = (1130, 950, 1290, 975)
BBOX_R4_GOLD = (1130, 995, 1290, 1020)
BBOX_R5_GOLD = (1130, 1040, 1290, 1065)

BBOX_B1_BOUNTY = (820, 847, 860, 870)
BBOX_B2_BOUNTY = (820, 893, 860, 915)
BBOX_B3_BOUNTY = (820, 937, 860, 957)
BBOX_B4_BOUNTY = (820, 981, 860, 1001)
BBOX_B5_BOUNTY = (820, 1025, 860, 1045)
BBOX_R1_BOUNTY = (1076, 847, 1118, 870)
BBOX_R2_BOUNTY = (1076, 893, 1118, 915)
BBOX_R3_BOUNTY = (1076, 937, 1118, 957)
BBOX_R4_BOUNTY = (1076, 981, 1118, 1001)
BBOX_R5_BOUNTY = (1076, 1025, 1118, 1045)
BBOX_R5_BOUNTY = (1076, 1025, 1118, 1045)

BBOX_B1_CREEP = (877, 860, 915, 885)
BBOX_B2_CREEP = (877, 905, 915, 930)
BBOX_B3_CREEP = (877, 950, 915, 975)
BBOX_B4_CREEP = (877, 995, 915, 1020)
BBOX_B5_CREEP = (877, 1040, 915, 1065)
BBOX_R1_CREEP = (1018, 860, 1055, 885)
BBOX_R2_CREEP = (1018, 905, 1055, 930)
BBOX_R3_CREEP = (1018, 950, 1055, 975)
BBOX_R4_CREEP = (1018, 995, 1055, 1020)
BBOX_R5_CREEP = (1018, 1040, 1055, 1065)


bbox_team_gold = [
    {"bbox":BBOX_BLUE_TEAM_GOLD, "mask_type":"blue", "message":"Blue Team Gold detected: "},
    {"bbox":BBOX_RED_TEAM_GOLD, "mask_type":"red", "message":"Red Team Gold detected: "},
]

bbox_blue_player_gold = [
    {"bbox":BBOX_B1_GOLD, "mask_type":"white", "message":"B1 Gold detected: "},
    {"bbox":BBOX_B2_GOLD, "mask_type":"white", "message":"B2 Gold detected: "},
    {"bbox":BBOX_B3_GOLD, "mask_type":"white", "message":"B3 Gold detected: "},
    {"bbox":BBOX_B4_GOLD, "mask_type":"white", "message":"B4 Gold detected: "},
    {"bbox":BBOX_B5_GOLD, "mask_type":"white", "message":"B5 Gold detected: "},
]

bbox_red_player_gold = [
    {"bbox":BBOX_R1_GOLD, "mask_type":"white", "message":"R1 Gold detected: "},
    {"bbox":BBOX_R2_GOLD, "mask_type":"white", "message":"R2 Gold detected: "},
    {"bbox":BBOX_R3_GOLD, "mask_type":"white", "message":"R3 Gold detected: "},
    {"bbox":BBOX_R4_GOLD, "mask_type":"white", "message":"R4 Gold detected: "},
    {"bbox":BBOX_R5_GOLD, "mask_type":"white", "message":"R5 Gold detected: "},
]

bbox_blue_player_bounty = [
    {"bbox":BBOX_B1_BOUNTY, "mask_type":"gold", "message":"B1 Bounty detected: "},
    {"bbox":BBOX_B2_BOUNTY, "mask_type":"gold", "message":"B2 Bounty detected: "},
    {"bbox":BBOX_B3_BOUNTY, "mask_type":"gold", "message":"B3 Bounty detected: "},
    {"bbox":BBOX_B4_BOUNTY, "mask_type":"gold", "message":"B4 Bounty detected: "},
    {"bbox":BBOX_B5_BOUNTY, "mask_type":"gold", "message":"B5 Bounty detected: "},
]

bbox_red_player_bounty = [
    {"bbox":BBOX_R1_BOUNTY, "mask_type":"gold", "message":"R1 Bounty detected: "},
    {"bbox":BBOX_R2_BOUNTY, "mask_type":"gold", "message":"R2 Bounty detected: "},
    {"bbox":BBOX_R3_BOUNTY, "mask_type":"gold", "message":"R3 Bounty detected: "},
    {"bbox":BBOX_R4_BOUNTY, "mask_type":"gold", "message":"R4 Bounty detected: "},
    {"bbox":BBOX_R5_BOUNTY, "mask_type":"gold", "message":"R5 Bounty detected: "},
]

bbox_blue_player_creep = [
    {"bbox":BBOX_B1_CREEP, "mask_type":"gold", "message":"B1 CreepScore detected: "},
    {"bbox":BBOX_B2_CREEP, "mask_type":"gold", "message":"B2 CreepScore detected: "},
    {"bbox":BBOX_B3_CREEP, "mask_type":"gold", "message":"B3 CreepScore detected: "},
    {"bbox":BBOX_B4_CREEP, "mask_type":"gold", "message":"B4 CreepScore detected: "},
    {"bbox":BBOX_B5_CREEP, "mask_type":"gold", "message":"B5 CreepScore detected: "},
]

bbox_red_player_creep = [
    {"bbox":BBOX_R1_CREEP, "mask_type":"gold", "message":"R1 CreepScore detected: "},
    {"bbox":BBOX_R2_CREEP, "mask_type":"gold", "message":"R2 CreepScore detected: "},
    {"bbox":BBOX_R3_CREEP, "mask_type":"gold", "message":"R3 CreepScore detected: "},
    {"bbox":BBOX_R4_CREEP, "mask_type":"gold", "message":"R4 CreepScore detected: "},
    {"bbox":BBOX_R5_CREEP, "mask_type":"gold", "message":"R5 CreepScore detected: "},
]

def captureScreenArea(bbox, monitor_index=2):
    """
    Greift einen bestimmten Bereich (bbox) auf einem spezifischen Monitor ab.
    monitor_index: 1 für Hauptmonitor, 2 für den zweiten Monitor.
    bbox: Format (Start_X, Start_Y, End_X, End_Y)
    """
    with mss.mss() as sct:
        # Prüfen, ob der gewünschte Monitor existiert
        if monitor_index < len(sct.monitors):
            monitor = sct.monitors[monitor_index]
            
            # Umrechnen des Tuple-Formats (x1, y1, x2, y2) in das mss-Dictionary-Format
            # mss benötigt relative Koordinaten zum jeweiligen Monitor
            x1, y1, x2, y2 = bbox
            width = x2 - x1
            height = y2 - y1
            
            cfg = {
                "top": monitor["top"] + y1,
                "left": monitor["left"] + x1,
                "width": width,
                "height": height
            }
            
            screenshot = sct.grab(cfg)
            frame = np.array(screenshot)
            
            # mss liefert BGRA -> Konvertierung zu klassischem BGR für OpenCV
            return cv2.cvtColor(frame, cv2.COLOR_BGRA2BGR)
        else:
            print(font=f"Monitor {monitor_index} nicht gefunden! Nutze Hauptbildschirm.")
            # Fallback auf den Hauptmonitor (Index 1), falls Monitor 2 fehlt
            return None

def createMask(image, mask_type):
    hsv = cv2.cvtColor(image, cv2.COLOR_BGR2HSV)
    match mask_type:
        case "blue":
            lower_blue = np.array([90, 100, 100])
            upper_blue = np.array([120, 255, 255])
            mask = cv2.inRange(hsv, lower_blue, upper_blue)
        case "red":
            lower_red1 = np.array([0, 150, 80])
            upper_red1 = np.array([10, 255, 255])

            lower_red2 = np.array([170, 150, 80])
            upper_red2 = np.array([180, 255, 255])

            mask1 = cv2.inRange(hsv, lower_red1, upper_red1)
            mask2 = cv2.inRange(hsv, lower_red2, upper_red2)
            mask = cv2.bitwise_or(mask1, mask2)
        case "white":
            lower_white = np.array([0,0,180])
            upper_white = np.array([180, 50, 255])
            mask = cv2.inRange(hsv, lower_white, upper_white)
        case "gold":
            lower_gold = np.array([15, 60, 110])
            upper_gold = np.array([30, 200, 255])
            mask = cv2.inRange(hsv, lower_gold, upper_gold)
    return mask

def preprocessImage(image, mask_type):
    if image is None:
        return None
    resized = cv2.resize(image, None, fx=3, fy=3, interpolation=cv2.INTER_CUBIC)
    processed = cv2.GaussianBlur(createMask(resized, mask_type), (3,3), 0)
    return processed

def getNumbersFromScreen(bbox, mask_type):
    # Hier sagen wir der Funktion, dass sie von Monitor 2 lesen soll
    cropped = captureScreenArea(bbox, monitor_index=3)
    
    if cropped is None:
        return "Error (No Monitor)"
        
    processed = preprocessImage(cropped, mask_type)
    cv2.imshow("Debug", processed)
    cv2.waitKey(1)
    customConfig = r'--oem 3 --psm 7 -c tessedit_char_whitelist=0123456789.kK()'
    text = pytesseract.image_to_string(processed, config=customConfig)
    return text.strip()

def parseGold(gold_value):
    matches = re.findall(r'\d+', gold_value)
    numbers = [int(num) for num in matches]
    return numbers

def getBluePlayerGold():
    data = {}
    for i, e in enumerate(bbox_blue_player_gold):
        value = getNumbersFromScreen(e["bbox"], e["mask_type"])
        gold_values = parseGold(value)
        key = f"b{i+1}"
        data[key] = gold_values[1]
    return data

def getRedPlayerGold():
    data = {}
    for i, e in enumerate(bbox_red_player_gold):
        value = getNumbersFromScreen(e["bbox"], e["mask_type"])
        gold_values = parseGold(value)
        key = f"r{i+1}"
        data[key] = gold_values[1]
    return data

def getBluePlayerCreep():
    data = {}
    for i, e in enumerate(bbox_blue_player_creep):
        value = getNumbersFromScreen(e["bbox"], e["mask_type"])
        key = f"b{i+1}"
        data[key] = value
    return data

def getRedPlayerCreep():
    data = {}
    for i, e in enumerate(bbox_red_player_creep):
        value = getNumbersFromScreen(e["bbox"], e["mask_type"])
        key = f"r{i+1}"
        data[key] = value
    return data

def getBluePlayerBounty():
    data = {}
    for i, e in enumerate(bbox_blue_player_bounty):
        value = getNumbersFromScreen(e["bbox"], e["mask_type"])
        key = f"b{i+1}"
        data[key] = value
    return data

def getRedPlayerBounty():
    data = {}
    for i, e in enumerate(bbox_red_player_bounty):
        value = getNumbersFromScreen(e["bbox"], e["mask_type"])
        key = f"r{i+1}"
        data[key] = value
    return data

def calcTotalTeamGold(blue_team_gold, red_team_gold):
    blue_team_total = 0
    for e in blue_team_gold.values():
        blue_team_total += e
    red_team_total = 0
    for e in red_team_gold.values():
        red_team_total += e
    data = {"blue_team":blue_team_total, "red_team":red_team_total}
    return data

def getData():
    data = {}

    blue_player_gold = getBluePlayerGold()
    red_player_gold = getRedPlayerGold()
    data["gold"] = {"blue":blue_player_gold, "red":red_player_gold}

    blue_player_creep = getBluePlayerCreep()
    red_player_creep = getRedPlayerCreep()
    data["creep"] = {"blue":blue_player_creep, "red":red_player_creep}

    blue_player_bounty = getBluePlayerBounty()
    red_player_bounty = getRedPlayerBounty()
    data["bounty"] = {"blue":blue_player_bounty, "red":red_player_bounty}

    total_team_gold = calcTotalTeamGold(blue_player_gold, red_player_gold)
    data["team_gold"] = total_team_gold

    return data


if __name__ == "__main__":
    print("Script active. Scanning MONITOR 2 for League of Legends...")
    

    while True:
        data = getData()
        print(f"data: {data}")
        time.sleep(10)
    # pydirectinput.press('o') # Falls du das Scoreboard umschalten willst