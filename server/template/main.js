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
        start: offset,
        end: pageSize(),
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
    const response = await fetch(url, {method: "DELETE"});
    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
    }
}

async function deleteJobRequest(queue, id) {
    const url = `/api/v1/queues/${queue}/jobs/${id}`;
    const response = await fetch(url, {method: "DELETE"});
    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
    }
}
async function retryJobRequest(queue, id) {
    const url = `/api/v1/queues/${queue}/jobs/${id}`;
    const response = await fetch(url, {method: "POST"});
    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
    }
}

async function getApi(path, filter) {
    const url = `/api/v1/${path}?${query(filter)}`;
    const signal = abortController.signal;
    const response = await fetch(url, {signal});
    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
    }

    return await response.json()
}

function setPageSize() {
    localStorage.setItem('pageSize', parseInt(document.getElementById("pageSize").value));
}

function pageSize() {
    return localStorage.getItem('pageSize');
}

function clearQueue(queue) {
    clearQueueRequest(queue);
    loadQueues();
}

function deleteJob(queue, id) {
    console.log("Tried to delete: " + id)
    deleteJobRequest(queue, id)
}

function retryJob(queue, id) {
    console.log("Tried to retry: " + id)
    retryJobRequest(queue, id)
}

/**
 * Worker methods
 */
function getWorkerRow(key, item) {
    let entry = ""
    if (item["entry"]["class"] !== "") {
        entry = JSON.stringify(item["entry"])
    }
    return `<tr><td>${key}</td><td>${item["host"]}</td><td>${item["socket"]}</td><td>${entry}</td></tr>`;
}
/**
 * Queue methods
 */

function getQueueRow(item) {
    return `<tr><td>${item["name"]}</td><td>${item["job_count"]}</td><td><a href="/?queue=${item["id"]}" class="btn btn-primary">Jobs</a> <button onclick="clearQueue('${item["id"]}')" class="btn btn-danger">Clear</button></td></tr>`;
}

/**
 * Job methods
 */

function getJobsHeader(failed) {
    let additionalHtml = ''
    if (failed) {
        additionalHtml = `<th scope="col">Exception</th><th scope="col">Failed At</th>`;
    }

    return `<tr><th scope="col"><input type="checkbox" class="form-check-inline" id="check-all"/></th><th scope="col">Class</th><th scope="col">Queued At</th>${additionalHtml}<th></th></tr>`;
}

function getJobClassSelect(items, filter) {
    let html = `<option ${filter.class ? '' : 'selected'} value="">-- Select class --</option>`;
    for (let key in items) {
        let count = items[key];
        html += `<option value='${key}' ${(filter.class === key) ? 'selected' : ''}>${key} (${count} items)</option>`;
    }

    return html
}

function getJobExceptionSelect(items, filter) {
    let html = `<option ${filter.class ? '' : 'selected'} value="">-- Select exception --</option>`;
    for (let key in items) {
        let count = items[key];
        html += `<option value='${key}' ${(filter.exception === key) ? 'selected' : ''}>${key} (${count} items)</option>`;
    }

    return html
}

function getJobRow(item, failed) {
    let job = failed ? item.payload : item
    let date = new Date(job.queue_time * 1000);

    let additionalHtml = ''
    if (failed) {
        additionalHtml = `<td>${item.exception}: ${item.error}</td><td>${item.failed_at}</td>`
    }

    return `<tr><td><input type="checkbox" class="form-check-inline" id="check-${job.id}"/></td><td>${job.class}</td><td>${date.toISOString()}</td>${additionalHtml}<td><button class="btn btn-outline-info" type="button" data-bs-toggle="modal" data-bs-target="#detailModal-${job.id}">Details</button></td></tr>`;
}

function getJobModal(item, failed) {
    let job = failed ? item.payload : item
    let date = new Date(job.queue_time * 1000);

    let additionalModalFields = ''
    if (failed) {
        additionalModalFields = `<label for="queue-${job.id}" class="form-label">Queue</label>
                                             <input type="text" id="queue-${job.id}" class="form-control" readonly value="${item.queue}"/>
                                             <label for="worker-${job.id}" class="form-label">Worker</label>
                                             <input type="text" id="worker-${job.id}" class="form-control" readonly value="${item.worker}"/>
                                             <label for="failed-${job.id}" class="form-label">Failed</label>
                                             <input type="text" id="failed-${job.id}" class="form-control" readonly value="${item.failed_at}"/>
                                             <label for="exception-${job.id}" class="form-label">Exception</label>
                                             <input type="text" id="exception-${job.id}" class="form-control" readonly value="${item.exception}"/>
                                             <label for="message-${job.id}" class="form-label">Message</label>
                                             <input type="text" id="message-${job.id}" class="form-control" readonly value="${item.error}"/>
                                             <label for="backtrace-${job.id}" class="form-label">Backtrace</label>
                                             <textarea readonly id="backtrace-${job.id}" class="form-control" style="height: 250px;">${item.backtrace.join('\n')}</textarea>`;
    }

    return `<div class="modal fade" id="detailModal-${job.id}" tabindex="-1" aria-labelledby="modalLabel-${job.id}" aria-hidden="true">
              <div class="modal-dialog modal-xl">
                <div class="modal-content">
                  <div class="modal-header">
                    <h1 class="modal-title fs-5" id="modalLabel-${job.id}">Job details</h1>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                  </div>
                  <div class="modal-body">
                    <label for="class-${job.id}" class="form-label">Class</label>
                    <input type="text" id="class-${job.id}" class="form-control" readonly value="${job.class}"/>
                    <label for="id-${job.id}" class="form-label">Id</label>
                    <input type="text" id="id-${job.id}" class="form-control" readonly value="${job.id}"/>
                    <label for="queued-${job.id}" class="form-label">Queued</label>
                    <input type="text" id="queued-${job.id}" class="form-control" readonly value="${date.toISOString()}"/>
                    <label for="args-${job.id}" class="form-label">Arguments</label>
                    <textarea readonly id="args-${job.id}" class="form-control" style="height: 250px;">${JSON.stringify(job.args, null, 2)}</textarea>
                    ${additionalModalFields}
                  </div>
                  <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                    <button type="button" class="btn btn-danger" onclick="deleteJob('failed', '${job.id}')">Delete</button>
                    <button type="button" class="btn btn-warning" onclick="retryJob('failed', '${job.id}')">Retry</button>
                  </div>
                </div>
              </div>
            </div>`;
}