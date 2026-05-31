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

    if (blue.gold) {
        document.querySelector('.blue-side .team-gold').textContent = formatGold(blue.gold);
    } else document.querySelector('.blue-side .team-gold').textContent = "0.0K";
    document.querySelector('.blue-kills').textContent = blue.score;

    if (red.gold) {
        document.querySelector('.red-side .team-gold').textContent = formatGold(red.gold);
    } else document.querySelector('.red-side .team-gold').textContent = "0.0K";
    document.querySelector('.red-kills').textContent = red.score;

    const goldDiff = blue.gold - red.gold;
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

        if (bluePlayerData && redPlayerData) {
            const diffEl = row.querySelector('.lane-gold-diff');
            const diffAmountEl = row.querySelector('.lane-gold-diff .diff-amount');

            if (diffEl && diffAmountEl) {
                const laneDiff = bluePlayerData.playerTotalGold - redPlayerData.playerTotalGold;
                diffAmountEl.textContent = formatGold(Math.abs(laneDiff));

                if (laneDiff > 0) {
                    diffEl.classList.remove('red-lead');
                    diffEl.classList.add('blue-lead');
                } else if (laneDiff < 0) {
                    diffEl.classList.remove('blue-lead');
                    diffEl.classList.add('red-lead');
                } else {
                    diffEl.classList.remove('blue-lead');
                    diffEl.classList.remove('red-lead');
                    diffEl.classList.add('no-lead');
                    diffAmountEl.textContent = "0.0K";
                }
            }
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
        const csEl = cardEl.querySelector('.player-cs');
        if (csEl) csEl.textContent = p.scores.creepScore;

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
            wardSlot.classList.remove("empty-slot")
            wardSlot.innerHTML = `
                <img src="${wardItem.icon}" alt="${wardItem.id}">
                <span class="ward-count">${Math.floor(p.scores.wardScore)}</span>
            `;
        } else {
            wardSlot.innerHTML = `<span class="ward-count">${Math.floor(p.scores.wardScore)}</span>`;
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