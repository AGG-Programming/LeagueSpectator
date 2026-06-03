const socket = new WebSocket('ws://localhost:8080/ws');

socket.onopen = () => {
    console.log('WebSocket-Connection to the Spactator Server established.');
};

socket.onmessage = (event) => {
    try {
        const data = JSON.parse(event.data);

        updateGameMeta(data);
        updateScoreDisplay(data);
        updateTimers(data.timers);
        updatePlayerScoreboard(data);
    }   catch (error) {
        console.error('Error parsing the WebSocket data:', error);
    }
};

socket.onclose = () => {
    console.warn('WebSocket-Connection closed. Try to reconnect in 5 seconds...');
    setTimeout(() => {
        window.location.reload();
    }, 5000);
};

socket.onerror = (error) => {
    console.error('WebSocket error:', error);
};

const scriptSocket = new WebSocket('ws://localhost:8765/ws')

scriptSocket.onopen = () => {
    console.log('WebSocket-Connection for the P-Script to the Spactator Server established.');
};

scriptSocket.onmessage = (event) => {
    try {
        const data = JSON.parse(event.data);

        updateTeamGold(data);
        updatePlayerGoldDiff(data);
        updatePlayerCs(data);
        updatePlayerBounty(data);

    }   catch (error) {
        console.error('Error parsing the WebSocket for P-Script data:', error);
    }
};

scriptSocket.onclose = () => {
    console.warn('WebSocket-Connection for P-Script closed. Try to reconnect in 5 seconds...');
    setTimeout(() => {
        window.location.reload();
    }, 5000);
};

scriptSocket.onerror = (error) => {
    console.error('WebSocket for P-Script error:', error);
};

function formatGameTime(seconds) {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
}

function formatGold(goldAmount) {
    return `${(goldAmount / 1000).toFixed(1)}K`;
}

function updateGameMeta(data) {
    if (!data || !data.gameTime) return;

    const timeEl = document.querySelector('.game-time');
    if (timeEl) {
        timeEl.textContent = formatGameTime(data.gameTime);
    }
}

function updateScoreDisplay(data) {
    const blue = data.blueTeam;
    const red = data.redTeam;

    if (!blue || !red) return;

    document.querySelector('.blue-kills').textContent = blue.score;

    document.querySelector('.red-kills').textContent = red.score;

    if (blue.objectives) {
        blue.objectives.forEach(obj => {
            if (obj.key === 'turrets') document.querySelector('.blue-towers .obj-count').textContent = obj.kills;
            if (obj.key === 'grubs') document.querySelector('.blue-grubs .obj-count').textContent = obj.kills;

            if (obj.key === 'dragon') {
                const blueDrakesContainer = document.querySelector('.blue-drakes');
                if (blueDrakesContainer) {
                    blueDrakesContainer.innerHTML = `<img src="${obj.icon}" class="drake-icon" alt="${obj.name}">`;
                }
            }
        });
    }

    if (red.objectives) {
        red.objectives.forEach(obj => {
            if (obj.key === 'turrets') document.querySelector('.red-towers .obj-count').textContent = obj.kills;
            if (obj.key === 'grubs') document.querySelector('.red-grubs .obj-count').textContent = obj.kills;

            if (obj.key === 'dragon') {
                const redDrakesContainer = document.querySelector('.red-drakes');
                if (redDrakesContainer) {
                    redDrakesContainer.innerHTML = `<img src="${obj.icon}" class="drake-icon" alt="${obj.name}">`;
                }
            }
        });
    }
}


function updateTimers(timersArray) {
    if (!timersArray || !Array.isArray(timersArray)) return;

    timersArray.forEach(timer => {
        if (timer.type === 'dragon') {
            const drakeTimeEl = document.querySelector('#drake-timer-box .spawn-time');
            if (drakeTimeEl) {
                drakeTimeEl.textContent = timer.alive ? "LIVE" : formatGameTime(timer.SpawnTime);
                if (timer.alive) {
                    drakeTimeEl.classList.add('live');
                } else {
                    drakeTimeEl.classList.remove('live');
                }
            }
        }

        if (timer.type === 'baron') {
            const baronTimeEl = document.querySelector('#void-timer-box .spawn-time');
            if (baronTimeEl) {
                baronTimeEl.textContent = timer.alive ? "LIVE" : formatGameTime(timer.SpawnTime);
                if (timer.alive) {
                    baronTimeEl.classList.add('live');
                } else {
                    baronTimeEl.classList.remove('live');
                }
            }
        }
    });
}


function updatePlayerScoreboard(data) {
    const bluePlayers = data.blueTeam.players;
    const redPlayers = data.redTeam.players;

    if (!bluePlayers || !redPlayers) return;

    const rows = document.querySelectorAll('#player-scoreboard .player-row');

    for (let i = 0; i < 5; i++) {
        const row = rows[i];
        if (!row) break;

        const bluePlayerData = bluePlayers[i];
        const redPlayerData = redPlayers[i];

        const blueCard = row.querySelector('.blue-player');
        if (blueCard && bluePlayerData) {
            updatePlayerCard(blueCard, bluePlayerData);
        }

        const redCard = row.querySelector('.red-player');
        if (redCard && redPlayerData) {
            updatePlayerCard(redCard, redPlayerData);
        }
    }
}


function updatePlayerCard(cardEl, p) {
    const champImg = cardEl.querySelector('.champ-img');
    if (champImg) champImg.src = `${p.icon}`;

    const lvlEl = cardEl.querySelector('.champ-level');
    if (lvlEl) lvlEl.textContent = p.level;

    const nameEl = cardEl.querySelector('.player-name');
    if (nameEl) nameEl.textContent = p.riotId ? p.riotId.split('#')[0] : '';

    if (p.scores) {
        const kdaEl = cardEl.querySelector('.player-kda');
        if (kdaEl) kdaEl.textContent = `${p.scores.kills}/${p.scores.deaths}/${p.scores.assists}`;
    }

    const avatarBox = cardEl.querySelector('.champ-avatar-box');
    if (avatarBox) {
        const deathOverlay = cardEl.querySelector('.death-timer-overlay');
        const hpBar = cardEl.querySelector('.hp-bar');
        const manaBar = cardEl.querySelector('.mana-bar');
        const bountyTag = cardEl.querySelector('.bounty-tag');

        if (p.isDead) {
            avatarBox.classList.add('dead');
            if (hpBar) {
                hpBar.style.width = '0%';
                hpBar.classList.add('dead-bar');
            }
            if (manaBar) {
                manaBar.style.width = '0%'
            }
            if (deathOverlay) {
                deathOverlay.classList.remove('hidden');
                deathOverlay.innerHTML = `<span>${Math.floor(p.respawnTimer)}</span>`;
            }
        } else {
            avatarBox.classList.remove('dead');
            if (hpBar) {
                hpBar.style.width = '100%';
                hpBar.classList.remove('dead-bar')
            }
            if (manaBar) {
                manaBar.style.width = '100%';
            }
            if (deathOverlay) deathOverlay.classList.add('hidden');
        }
    }

    if (p.runes) {
        const runeImages = cardEl.querySelectorAll('.runes-box img');
        if (runeImages.length >= 2) {
            if (p.runes.keystone) runeImages[0].src = p.runes.keystone.icon;
            if (p.runes.secondary) runeImages[1].src = p.runes.secondary.icon;
        }
    }

    const summSlots = cardEl.querySelectorAll('.summs-container .summ-slot');
    if (summSlots.length >= 2) {
        if (p.spells.length >= 2) {
            summSlots[0].innerHTML = `<img src="${p.spells[0].icon}" alt="${p.spells[0].displayName}">`;
            summSlots[1].innerHTML = `<img src="${p.spells[1].icon}" alt="${p.spells[1].displayName}">`;
        }
    }

    const itemSlots = cardEl.querySelectorAll('.items-line .item-slot');
    itemSlots.forEach(slot => {
        slot.className = 'item-slot empty';
        slot.innerHTML = '';
    });

    if (p.items && Array.isArray(p.items)) {
        p.items.forEach(item => {
            const slotIdx = item.slot;
            if (slotIdx >= 6) return;
            const slot = itemSlots[slotIdx];
            if (!slot) return;

            slot.className = 'item-slot';
            slot.innerHTML = `<img src="${item.icon}" alt="${item.id}">`;

            if (item.consumable) {
                slot.innerHTML += `<span class="item-count">${item.count}</span>`
            }
        });
    }

    const wardSlot = cardEl.querySelector('.ward-slot');
    if (wardSlot) {
        const wardItem = p.items.find(item => item.slot === 6);
        if (wardItem) {
            wardSlot.classList.remove("empty-slot");
            wardSlot.innerHTML = `
                <img src="${wardItem.icon}" alt="${wardItem.id}">
                <span class="ward-count">${Math.floor(p.scores.wardScore)}</span>
            `;
        } else {
            wardSlot.innerHTML = `<span class="ward-count">${Math.floor(p.scores.wardScore)}</span>`;
        }
    }

    const ultSlot = cardEl.querySelector('.ult-slot');
    if (ultSlot) {
        if (p.ultIcon) {
            ultSlot.innerHTML = `<img src="${p.ultIcon}" alt="U">`;
            if (p.level >= 6) {
                ultSlot.classList.remove("unlearned");
                ultSlot.classList.add("ready");
            } else {
                ultSlot.classList.remove("ready");
                ultSlot.classList.add("unlearned");
            }
        }
    }
}



//-----------------------------WebSocket P-Script------------------------------------#

function updateTeamGold(data) {
    const blue_gold = data.team_gold.blue_team;
    const red_gold = data.team_gold.red_team;

    if (blue_gold) {
        document.querySelector('.blue-side .team-gold').textContent = formatGold(blue_gold);
    } else document.querySelector('.blue-side .team-gold').textContent = "0.0K";

    if (red_gold) {
        document.querySelector('.red-side .team-gold').textContent = formatGold(red_gold);
    } else document.querySelector('.red-side .team-gold').textContent = "0.0K";

    const goldDiff = blue_gold - red_gold;
    const blueLeadEl = document.querySelector('.blue-side .gold-lead');
    const redLeadEl = document.querySelector('.red-side .gold-lead');

    if (goldDiff > 0) {
        if (blueLeadEl) {
            blueLeadEl.classList.remove('hidden');
            blueLeadEl.classList.add('blue-side-lead');
            blueLeadEl.textContent = `+${formatGold(goldDiff)}`;
        }
        if (redLeadEl) redLeadEl.classList.add('hidden');
    } else if (goldDiff < 0) {
        if (redLeadEl) {
            redLeadEl.classList.remove('hidden');
            redLeadEl.classList.add('red-side-lead');
            redLeadEl.textContent = `+${formatGold(Math.abs(goldDiff))}`;
        }
        if (blueLeadEl) blueLeadEl.classList.add('hidden');
    } else {
        if (blueLeadEl) blueLeadEl.classList.add('hidden');
        if (redLeadEl) redLeadEl.classList.add('hidden');
    }
}

function updatePlayerGoldDiff(data) {
    const bluePlayers = data.gold.blue;
    const redPlayers = data.gold.red;

    if (!bluePlayers || !redPlayers) return;

    const rows = document.querySelectorAll('#player-scoreboard .player-row');

    for (let i = 0; i < 5; i++) {
        const row = rows[i];
        if (!row) break;

        const bPlayer = `b${i+1}`;
        const rPlayer = `r${i+1}`;
        const blueGold = bluePlayers[bPlayer];
        const redGold = redPlayers[rPlayer];
        console.log("blueGold: ",blueGold);
        console.log("redGold: ",redGold);
        const bluePlayerGold = parseInt(blueGold, 10);
        const redPlayerGold = parseInt(redGold, 10);
        console.log("blueGold: ",bluePlayerGold);
        console.log("redGold: ",redPlayerGold);

        const goldDiff = bluePlayerGold - redPlayerGold;
        const LeadEl = row.querySelector('.lane-gold-diff');
        console.log("gold diff: ", goldDiff);

        if (LeadEl) {
            const amountEl = LeadEl.querySelector('.diff-amount');
            
            // 1. Erstmal alle Lead-Klassen komplett säubern
            LeadEl.classList.remove('blue-lead', 'red-lead', 'no-lead');

            if (goldDiff > 0) {
                LeadEl.classList.add('blue-lead');
                if (amountEl) amountEl.textContent = `+${formatGold(goldDiff)}`;
            } else if (goldDiff < 0) {
                LeadEl.classList.add('red-lead');
                if (amountEl) amountEl.textContent = `+${formatGold(Math.abs(goldDiff))}`;
            } else {
                LeadEl.classList.add('no-lead');
                if (amountEl) amountEl.textContent = '0.0K';
            }
        }
    }
}

function updatePlayerCs(data) {
    const bluePlayers = data.creep.blue;
    const redPlayers = data.creep.red;

    if (!bluePlayers || !redPlayers) return;

    const rows = document.querySelectorAll('#player-scoreboard .player-row');

    for (let i = 0; i < 5; i++) {
        const row = rows[i];
        if (!row) break;

        const bPlayer = `b${i+1}`;
        const rPlayer = `r${i+1}`;
        const bluePlayerData = bluePlayers[bPlayer];
        const redPlayerData = redPlayers[rPlayer];

        const blueCard = row.querySelector('.blue-player');
        if (blueCard && bluePlayerData) {
            const csEl = blueCard.querySelector('.player-cs');
            if (csEl) csEl.textContent = bluePlayerData;
        }

        const redCard = row.querySelector('.red-player');
        if (redCard && redPlayerData) {
            const csEl = redCard.querySelector('.player-cs');
            if (csEl) csEl.textContent = redPlayerData;
        }
    }
}

function updatePlayerBounty(data) {
    const bluePlayers = data.bounty.blue;
    const redPlayers = data.bounty.red;

    if (!bluePlayers || !redPlayers) return;

    const rows = document.querySelectorAll('#player-scoreboard .player-row');

    for (let i = 0; i < 5; i++) {
        const row = rows[i];
        if (!row) break;

        const bPlayer = `b${i+1}`;
        const rPlayer = `r${i+1}`;
        const bluePlayerBounty = bluePlayers[bPlayer];
        const redPlayerBounty = redPlayers[rPlayer];

        const blueCard = row.querySelector('.blue-player');
        if (blueCard) {
            const bountyEl = blueCard.querySelector('.bounty-tag');
            if (bountyEl) {
                if (bluePlayerBounty > 0) {
                    bountyEl.classList.remove('hidden');
                    bountyEl.textContent = bluePlayerBounty;
                } else {
                    bountyEl.classList.add('hidden');
                    bountyEl.textContent = "";
                }
            }
        }

        const redCard = row.querySelector('.red-player');
        if (redCard) {
            const bountyEl = redCard.querySelector('.bounty-tag');
            if (bountyEl) {
                if (redPlayerBounty > 0) {
                    bountyEl.classList.remove('hidden');
                    bountyEl.textContent = redPlayerBounty;
                } else {
                    bountyEl.classList.add('hidden');
                    bountyEl.textContent = "";
                }
            }
        }
    }
}


//completed in div when completed

//                    <div class="summ-slot cd">
//                        <img src="assets/images/summs/cloud-drake.png" alt="F">
//                        <div class="cooldown-overlay"></div>

//                        <div class="ult-slot ready">
//                           <img src="assets/images/champs/cloud-drake.png" alt="U">
//                        </div>

//            <span class="spawn-time live">LIVE</span>

//                <div class="drakes-list blue-drakes">
//                    <img src="assets/images/ui/normal-drake/32px-Infernal_Dragon_Soul_buff - normal.png" class="drake-icon" alt="I">
//                    <img src="assets/images/ui/normal-drake/32px-Ocean_Dragon_Soul_buff - normal.png" class="drake-icon" alt="O">
//                </div>

//                    <span class="gold-lead blue-side-lead">+1.4K</span>

//            <div class="lane-gold-diff blue-lead">
//                <span class="diff-amount">0.0K</span>

/*                    <div class="ward-slot">
                        <img src="assets/images/items/cloud-drake.png" alt="W">
                        <span class="ward-count">3</span>
                    </div>
                    <div class="quest-slot"><img src="assets/images/items/cloud-drake.png" alt="Q"></div>
                </div>

                <div class="items-line">
                    <div class="item-slot">
                        <img src="assets/images/items/cloud-drake.png" alt="Black Cleaver">
                        <span class="item-count">2</span>
                    </div>
                    <div class="item-slot cd-item">
                        <img src="assets/images/items/cloud-drake.png" alt="BOTRK">
                        <div class="item-cooldown-overlay"></div>
                    </div>
                    <div class="item-slot empty"></div>
                    <div class="item-slot empty"></div>
                    <div class="item-slot empty"></div>
                    <div class="item-slot empty"></div>
                </div>
*/

/*
                <div class="summs-container">
                    <div class="summ-slot">
                        <img src="assets/images/summs/cloud-drake.png" alt="F">
                    </div>
                    <div class="summ-slot">
                        <img src="assets/images/summs/cloud-drake.png" alt="T">
                    </div>
                </div>
*/

/*
                <div class="runes-box">
                    <img src="assets/images/runes/cloud-drake.png" class="main-rune" alt="R">
                    <img src="assets/images/runes/cloud-drake.png" class="sub-rune" alt="R">
                </div>
*/

/*
                <div class="champ-avatar-box">
                    <img src="assets/images/champs/cloud-drake.png" class="champ-img" alt="Champ">
                    <div class="champ-level">5</div>
                    <div class="death-timer-overlay hidden"></div>
                    <div class="bounty-tag hidden"></div>
                </div>
*/