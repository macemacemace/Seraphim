


function ago(amount: number, word:string) {
    if (amount === 1){
        return `1 ${word} ago`;
    }

    return `${amount} ${word}s ago`
}

export function timeAgo(dateText: string) {
    const miliseconds = Date.now() - new Date(dateText).getTime();
    const minutes = Math.floor(miliseconds / 60000);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);


    if(minutes < 1){
        return "Just now";
    }

    if(hours < 1){
        return ago(minutes, "minute");
    }

    if(days < 1){
        return ago(hours, "hour");
    }
    return ago(days, "day")
}
