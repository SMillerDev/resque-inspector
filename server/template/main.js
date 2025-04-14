function generateFilter() {
    let queue = document.getElementById("queues").value;
    let regex = document.getElementById("regex").value;
    let className = document.getElementById("classes").value;
    let exception = "";
    if (!document.getElementById("exceptions-wrapper").classList.contains("d-none")) {
        exception = document.getElementById("exceptions").value;
    }

    return {
        regex: regex,
        class: className,
        exception: exception,
        queue: queue,
        startDate: 0,
        endDate: Date.now(),
    }
}

function query(obj) {
    let str = [];
    for (const p in obj) {
        if (obj.hasOwnProperty(p) && obj[p] != undefined) {
            str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
        }
    }

    return str.join("&");
}

async function clearQueueRequest(queue) {
    const url = `/api/v1/queues/${queue}`;
    const response = await fetch(url, { method: "DELETE" });
    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
    }
}

async function getApi(path, filter, start = 0, offset = 25) {
    const url = `/api/v1/${path}?${query(filter)}&start=${start}&offset=${offset}`;
    const response = await fetch(url);
    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
    }

    return await response.json()
}

function onlyUnique(value, index, array) {
    return array.indexOf(value) === index;
}