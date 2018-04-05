export function timeRound(value){
    return Math.round(value * 100) / 100
}

export function currentTimeString() {
    var currentDate = new Date();
    return currentDate.getHours() + ":"
        + currentDate.getMinutes() + ":" + currentDate.getSeconds();
}