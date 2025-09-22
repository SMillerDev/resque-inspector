let abortController = new AbortController();
let offset = 0;
let loading = false;
let rowHtml = "";
let modalList = "";
if (localStorage.getItem("pageSize") === null) {
    setPageSize();
}
let classes = {};
let exceptions = {};




/**
 * Handle the infinite scroll
 * @return void
 */
const handleInfiniteScroll = () => {
    const endOfPage = window.innerHeight + window.scrollY >= document.body.offsetHeight;
    if (endOfPage) {
        let nextStart = Number(offset+pageSize());
        console.debug(`[Scroll] Loading ${nextStart} until ${nextStart + pageSize()}`)
        loadJobs(nextStart);
    }
};

let scrollListener = throttle(handleInfiniteScroll, 100);

/**
 * Load data onto the page
 * @param {string} page Page to load
 * @return void
 */
function loadData(page) {
    if (page === 'jobs') {
        window.addEventListener("scroll", scrollListener);
    } else {
        window.removeEventListener("scroll", scrollListener)
    }
    switch (page) {
        case 'queues':
            loadQueues();
            document.getElementById("stats").classList.add("d-none");
            break;
        case 'workers':
            loadWorkers()
            document.getElementById("stats").classList.add("d-none");
            break;
        case 'jobs':
            loadJobs(0)
            document.getElementById("stats").classList.remove("d-none");
            break;
        default:
            console.error(`No data to load for ${page}`)
            break;
    }
}

/**
 * Load the jobs onto the page
 * @param {number} start Starting job
 * @return void
 */
function loadJobs(start) {
    if (loading) {abortController.abort("Reloading data");}
    offset = start
    loading = true;
    if (offset < 1) {
        modalList = ''
        rowHtml = ''
    }

    document.getElementById("loading-spinner").classList.remove("d-none");
    const filter = generateFilter(offset, offset+pageSize());
    getApi(`queues/${filter.queue}/jobs`, filter).then(function (json) {
        console.debug(`[API] Loaded jobs ${offset} till ${offset + pageSize()}`);
        loading = false;
        if (json === null) {
            console.error("No job data!")
            return
        }
        if (filter.queue === 'failed') {
            document.getElementById("exceptions").classList.remove("d-none");
        } else {
            document.getElementById("exceptions").classList.add("d-none");
        }
        for (let index in json["classes"]) {
            if (classes.hasOwnProperty(index) === false) { classes[index] = 0; }
            classes[index] += json["classes"][index]
        }
        for (let index in json["exceptions"]) {
            if (exceptions.hasOwnProperty(index) === false) { exceptions[index] = 0; }
            exceptions[index] += json["exceptions"][index]
        }

        document.getElementById("classes").innerHTML = getJobClassSelect(classes, filter);
        document.getElementById("exceptions").innerHTML = getJobExceptionSelect(exceptions, filter);

        json["items"].forEach(item => {
            rowHtml += getJobRow(item, filter.queue === 'failed')
            modalList += getJobModal(filter.queue, item);
        });

        document.getElementById("loading-spinner").classList.add("d-none");
        document.getElementById("modal-list").innerHTML = modalList;
        document.getElementById("job-header").innerHTML = getJobsHeader(filter.queue === 'failed');
        document.getElementById("job-list").innerHTML = rowHtml;
        document.getElementById("total-count").innerHTML = json["total"];
        document.getElementById("loaded-count").innerHTML = offset < json["total"] ? String(offset + json["items"].length) : json["total"];
    }).catch(function (reason) {
        document.getElementById("loading-spinner").classList.add("d-none");
        console.error(reason);
    });
}

/**
 * Load the queues onto the page
 * @return void
 */
function loadQueues() {
    document.getElementById("loading-spinner").classList.remove("d-none");
    getApi("queues", {}).then(function (json) {
        if (json === null) {
            console.error("No Queue data!")
            return
        }
        let html = "";
        json["items"].forEach(item => {
            html += getQueueRow(item)
        });
        document.getElementById("loading-spinner").classList.add("d-none");
        document.getElementById("queue-list").innerHTML = html;
    }).catch(function (reason) {
        console.error(reason)
    });
}

/**
 * Load the workers onto the page.
 * @return void
 */
function loadWorkers() {
    document.getElementById("loading-spinner").classList.remove("d-none");
    getApi("workers", {}).then(function (json) {
        if (json === null) {
            console.error("No worker data!")
            return
        }
        let html = "";
        for (const key in json["items"]) {
            const list = json["items"][key];
            list.forEach(item => {
                html += getWorkerRow(key, item);
            });
        }
        document.getElementById("loading-spinner").classList.add("d-none");
        document.getElementById("worker-list").innerHTML = html;
    }).catch(function (reason) {
        console.error(reason)
    });
}

/**
 * Throttle in a timeframe
 * @param {function} func      The function to throttle
 * @param {number}   timeFrame The timeframe to throttle for
 * @return {(function(...[*]): void)|*}
 */
function throttle(func, timeFrame) {
    let lastTime = 0;
    return function (...args) {
        const now = new Date();
        if (now - lastTime >= timeFrame) {
            func(...args);
            lastTime = now;
        }
    };
}

/**
 * @typedef {Object} Filter
 * @property {string} regex     The regex to filter by
 * @property {string} class     The class name to filter by
 * @property {string} exception The exception to filter by
 * @property {string} queue     The queue to filter by
 * @property {number} startDate The start of the date range to filter by as epoch time
 * @property {number} endDate   The end of the date range to filter by as epoch time
 * @property {number} start     The start offset to filter by
 * @property {number} end       The limit of items to return
 */
/**
 * Generate a filter
 *
 * @param {number} start Start of the selection
 * @param {number} end   End of the selection
 *
 * @returns {Filter}
 */
function generateFilter(start, end) {
    let queue = document.getElementById("queues").value;
    let regex = document.getElementById("regex").value;
    let className = document.getElementById("classes").value;
    let exception = "";
    if (!document.getElementById("exceptions").classList.contains("d-none")) {
        exception = document.getElementById("exceptions").value;
    }

    return {
        regex: regex,
        class: className,
        exception: exception,
        queue: queue,
        startDate: 0,
        endDate: Date.now(),
        start: start,
        end: end,
    }
}

/**
 * Build a query from a filter
 * @param {object} filter
 * @returns {string}
 */
function query(filter) {
    let str = [];
    for (const p in filter) {
        if (filter.hasOwnProperty(p) && filter[p] !== undefined) {
            str.push(encodeURIComponent(p) + "=" + encodeURIComponent(filter[p]));
        }
    }

    return str.join("&");
}

/**
 * Toggle all checkboxes
 * @param {Node<HTMLElement>} source
 */
function toggleCheckboxes(source) {
    let checkboxes = document.getElementsByName('job-selector');
    for (let i = 0, n = checkboxes.length; i < n; i++) {
        checkboxes[i].checked = source.checked;
    }
}

/**
 * Show the edit banner
 * @param {Node<HTMLElement>} source
 */
function showEditBanner(source) {
    if (source.checked) {
        document.getElementById("edit-bar").classList.remove("d-none");
    } else {
        document.getElementById("edit-bar").classList.add("d-none");
    }
}

/**
 * Request a queue be cleared
 * @param {string} name
 * @returns {Promise<void>}
 */
async function clearQueueRequest(name) {
    const url = `/api/v1/queues/${name}`;
    const response = await fetch(url, {method: "DELETE"});
    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
    }
}

/**
 * Request a job be deleted
 * @param {string} queue
 * @param {string }id
 * @returns {Promise<void>}
 */
async function deleteJobRequest(queue, id) {
    const url = `/api/v1/queues/${queue}/jobs/${id}`;
    const response = await fetch(url, {method: "DELETE"});
    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
    }
}

/**
 * Request a job be retried
 * @param {string} queue
 * @param {string }id
 * @returns {Promise<void>}
 */
async function retryJobRequest(queue, id) {
    const url = `/api/v1/queues/${queue}/jobs/${id}`;
    const response = await fetch(url, {method: "POST"});
    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`);
    }
}

/**
 * Get an item from the API
 * @param {string} path
 * @param {object} filter
 * @returns {Promise<any>}
 */
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
    offset = 0;
    localStorage.setItem('pageSize', parseInt(document.getElementById("pageSize").value));
}

function pageSize() {
    return Number(localStorage.getItem('pageSize'));
}

/**
 * Clear queue
 * @param {string} name
 */
function clearQueue(name) {
    clearQueueRequest(name).then(() => loadQueues());
}

/**
 * Delete a job
 * @param {string}      queue Name of the queue
 * @param {string|null} id    ID of the job
 */
function deleteJob(queue, id) {
    let ids = []
    if (id === null) {
        let checkboxes = document.getElementsByName('job-selector');
        for (let i = 0, n = checkboxes.length; i < n; i++) {
            if (!checkboxes[i].checked) {
                continue;
            }
            ids.push(checkboxes[i].value)
        }
    } else {
        ids.push(id)
    }
    ids.forEach( (value) => deleteJobRequest(queue, value))
}

/**
 * Retry a job
 * @param {string}      queue Name of the queue
 * @param {string|null} id    ID of the job
 */
function retryJob(queue, id) {
    let ids = []
    if (id === null) {
        let checkboxes = document.getElementsByName('job-selector');
        for (let i = 0, n = checkboxes.length; i < n; i++) {
            if (!checkboxes[i].checked) {
                continue;
            }
            ids.push(checkboxes[i].value)
        }
    } else {
        ids.push(id)
    }
    ids.forEach( (value) => retryJobRequest(queue, value))
}

/* Worker methods */
/**
 * Get a row of workers
 * @param {string} key
 * @param {object} item
 * @returns {string}
 */
function getWorkerRow(key, item) {
    let entry = ""
    if (item["entry"]["class"] !== "") {
        entry = JSON.stringify(item["entry"])
    }
    return `<tr><td>${key}</td><td>${item.host}</td><td>${item.socket}</td><td>${entry}</td></tr>`;
}

/* Queue methods */

/**
 * Get a row of queues
 * @param {object} item
 * @returns {string}
 */
function getQueueRow(item) {
    return `<tr><td>${item.name}</td><td>${item["job_count"]}</td><td><a href="/jobs?queue=${item.id}" role="button">Jobs</a> <button onclick="clearQueue('${item.id}')" class="danger">Clear</button></td></tr>`;
}

/* Job methods */

/**
 * Get the correct header for jobs
 * @param {boolean} failed If the view shows failed items
 * @returns {string}
 */
function getJobsHeader(failed) {
    let additionalHtml = ''
    if (failed) {
        additionalHtml = `<th scope="col">Exception</th><th scope="col">Failed At</th>`;
    }

    return `<tr><th scope="col"><input type="checkbox" onclick="showEditBanner(this);toggleCheckboxes(this)" id="check-all"/></th><th scope="col">Class</th><th scope="col">Queued At</th>${additionalHtml}<th></th></tr>`;
}

/**
 * Get a selector for job classes
 * @param {object} items
 * @param {Filter} filter
 *
 * @returns {string}
 */
function getJobClassSelect(items, filter) {
    let selected = filter.class ? '' : 'selected'
    let html = `<option ${selected} disabled value="">-- Select Class --</option>`;
    for (let key in items) {
        let count = items[key];
        html += `<option value='${key}' ${(filter.class === key) ? 'selected' : ''}>${key} (${count} items)</option>`;
    }

    return html
}

/**
 * Get a selector for job exceptions
 * @param {object} items
 * @param {object} filter
 * @returns {string}
 */
function getJobExceptionSelect(items, filter) {
    let selected = filter.class ? '' : 'selected'
    let html = `<option ${selected} disabled value="">-- Select Exception --</option>`;
    for (let key in items) {
        let count = items[key];
        html += `<option value='${key}' ${(filter.exception === key) ? 'selected' : ''}>${key} (${count} items)</option>`;
    }

    return html
}

/**
 * Get a row for a job
 * @param {object} item
 * @param {boolean} failed
 * @returns {string}
 */
function getJobRow(item, failed) {
    let job = failed ? item.payload : item
    let date = new Date(job.queue_time * 1000);

    let additionalHtml = ''
    if (failed) {
        additionalHtml = `<td>${item.exception}: ${item.error}</td><td>${item.failed_at}</td>`
    }

    return `<tr>
                <td>
                    <input type="checkbox" name="job-selector" id="check-${job.id}" value="${job.id}" onclick="showEditBanner(this)"/>
                </td>
                <td>${job.class}</td>
                <td>${date.toISOString()}</td>
                ${additionalHtml}
                <td>
                    <button class="info" data-target="detailModal-${job.id}" onclick="toggleModal(event)">Details</button>
                </td>
            </tr>`;
}

/**
 * Get a modal for a job
 *
 * @param {string} queue Name of the queue
 * @param {object} item Item to parse
 *
 * @returns {string} The Job modal value
 */
function getJobModal(queue, item) {
    let failed = queue === 'failed'
    let job = failed ? item.payload : item
    let date = new Date(job.queue_time * 1000);

    let additionalModalFields = ''
    if (failed) {
        additionalModalFields = `<label for="queue-${job.id}" >Queue</label>
                                             <input type="text" id="queue-${job.id}" readonly value="${item.queue}"/>
                                             <label for="worker-${job.id}">Worker</label>
                                             <input type="text" id="worker-${job.id}" readonly value="${item.worker}"/>
                                             <label for="failed-${job.id}">Failed</label>
                                             <input type="text" id="failed-${job.id}" readonly value="${item.failed_at}"/>
                                             <label for="exception-${job.id}">Exception</label>
                                             <input type="text" id="exception-${job.id}" readonly value="${item.exception}"/>
                                             <label for="message-${job.id}">Message</label>
                                             <input type="text" id="message-${job.id}" readonly value="${item.error}"/>
                                             <label for="backtrace-${job.id}">Backtrace</label>
                                             <textarea readonly id="backtrace-${job.id}" style="height: 250px;">${item.backtrace.join('\n')}</textarea>`;
    }

    return `<dialog id="detailModal-${job.id}" aria-labelledby="modalLabel-${job.id}">
              <article>
                <header>
                    <h2 id="modalLabel-${job.id}">Job details</h2>
                </header>
                <form>
                    <label for="class-${job.id}">Class</label>
                    <input type="text" id="class-${job.id}" readonly value="${job.class}"/>
                    <label for="id-${job.id}">Id</label>
                    <input type="text" id="id-${job.id}" readonly value="${job.id}"/>
                    <label for="queued-${job.id}">Queued</label>
                    <input type="text" id="queued-${job.id}" readonly value="${date.toISOString()}"/>
                    <label for="args-${job.id}">Arguments</label>
                    <textarea readonly id="args-${job.id}" style="height: 250px;">${JSON.stringify(job.args, null, 2)}</textarea>
                    ${additionalModalFields}
                </form>
                <footer>
                    <button class="secondary" data-target="detailModal-${job.id}" onclick="toggleModal(event)">Close</button>
                    <button class="danger" onclick="deleteJob('${queue}', '${job.id}')">Delete</button>
                    <button class="warning" onclick="retryJob('${queue}', '${job.id}')">Retry</button>
                </footer>
              </article>
            </dialog>`;
}

/**
 * Pico CSS Modal handling
 */
const isOpenClass = "modal-is-open";
const openingClass = "modal-is-opening";
const closingClass = "modal-is-closing";
const scrollbarWidthCssVar = "--pico-scrollbar-width";
const animationDuration = 400; // ms
let visibleModal = null;

// Toggle modal
const toggleModal = (event) => {
    event.preventDefault();
    const modal = document.getElementById(event.currentTarget.dataset.target);
    if (!modal) return;
    modal && (modal.open ? closeModal(modal) : openModal(modal));
};

// Open modal
const openModal = (modal) => {
    const { documentElement: html } = document;
    const scrollbarWidth = getScrollbarWidth();
    if (scrollbarWidth) {
        html.style.setProperty(scrollbarWidthCssVar, `${scrollbarWidth}px`);
    }
    html.classList.add(isOpenClass, openingClass);
    setTimeout(() => {
        visibleModal = modal;
        html.classList.remove(openingClass);
    }, animationDuration);
    modal.showModal();
};

// Close modal
const closeModal = (modal) => {
    visibleModal = null;
    const { documentElement: html } = document;
    html.classList.add(closingClass);
    setTimeout(() => {
        html.classList.remove(closingClass, isOpenClass);
        html.style.removeProperty(scrollbarWidthCssVar);
        modal.close();
    }, animationDuration);
};

// Close with a click outside
document.addEventListener("click", (event) => {
    if (visibleModal === null) return;
    const modalContent = visibleModal.querySelector("article");
    const isClickInside = modalContent.contains(event.target);
    !isClickInside && closeModal(visibleModal);
});

// Close with Esc key
document.addEventListener("keydown", (event) => {
    if (event.key === "Escape" && visibleModal) {
        closeModal(visibleModal);
    }
});

// Get scrollbar width
const getScrollbarWidth = () => window.innerWidth - document.documentElement.clientWidth;

// Is scrollbar visible
const isScrollbarVisible = () => {
    return document.body.scrollHeight > screen.height;
};
