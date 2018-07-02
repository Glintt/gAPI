
let $ = require('jquery')

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
export function convertMillisToTime(millis){
    let delim = " ";
    let hours = Math.floor(millis / (1000 * 60 * 60) % 60);
    let minutes = Math.floor(millis / (1000 * 60) % 60);
    let seconds = Math.floor(millis / 1000 % 60);
    hours = hours < 10 ? '0' + hours : hours;
    minutes = minutes < 10 ? '0' + minutes : minutes;
    seconds = seconds < 10 ? '0' + seconds : seconds;
    return hours + 'h'+ delim + minutes + 'm' + delim + seconds + 's';
}

export function enableTooltip() {   
    $(document).ready(function() {
        $("body").tooltip({ selector: '[data-toggle=tooltip]', trigger : 'hover' });
    });
}