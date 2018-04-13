export function timeRound(value){
    return Math.round(value * 100) / 100
}

export function currentTimeString() {
    var currentDate = new Date();
    return currentDate.getHours() + ":"
        + currentDate.getMinutes() + ":" + currentDate.getSeconds();
}

export function getCookieByName(name){
    var pair = document.cookie.match(new RegExp(name + '=([^;]+)'));
    return pair ? pair[1] : null;
}