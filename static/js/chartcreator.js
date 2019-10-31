// Note about canvas sizes:
// https://www.chartjs.org/docs/latest/general/responsive.html#important-note

async function CreateMessagesPerMonthChart() {
    let resp = await fetch("/api/messages/permonth");
    resp = await resp.json();

    let data = {
        labels: [],
        datasets: [{
            data: [],
            borderColor: "rgb(75, 192, 192)",
            lineTension: 0.25
        }]
    };

    for (const year in resp) {
        for (const month in resp[year]) {
            data.labels.push(`${year}-${month}`);
            data.datasets[0].data.push(resp[year][month]);
        }
    }

    const ctx = document.getElementById("messagesPerMonthChart").getContext("2d");
    new Chart(ctx, {
        type: "line",
        data: data,
        options: {
            title: {
                display: true,
                text: "Messages Sent Per Month"
            },
            legend: {
                display: false
            },
            maintainAspectRatio: false,
        }
    });
}

async function CreateMessagesPerUserChart() {
    let resp = await fetch("/api/messages/peruser");
    resp = await resp.json();

    let data = {
        labels: [],
        datasets: [{
            data: [],
            backgroundColor: "rgb(75, 192, 192)"
        }]
    };

    for (const user in resp) {
        data.labels.push(user);
        data.datasets[0].data.push(resp[user]);
    }

    const ctx = document.getElementById("messagesPerUserChart").getContext("2d");
    new Chart(ctx, {
        type: "bar",
        data: data,
        options: {
            title: {
                display: true,
                text: "Messages Sent Per User"
            },
            legend: {
                display: false
            },
            maintainAspectRatio: false,
            scales: {
                yAxes: [{
                    ticks: {
                        beginAtZero: true
                    }
                }]
            }
        }
    });
}

async function CreateMessagesPerWeekdayChart() {
    let resp = await fetch("/api/messages/perweekday");
    resp = await resp.json();

    let data = {
        labels: [],
        datasets: [{
            data: [],
            borderColor: "rgb(75, 192, 192)"
        }]
    };

    for (const weekday in resp) {
        data.labels.push(weekday);
        data.datasets[0].data.push(resp[weekday]);
    }

    const ctx = document.getElementById("messagesPerWeekdayChart").getContext("2d");
    new Chart(ctx, {
        type: "radar",
        data: data,
        options: {
            title: {
                display: true,
                text: "Messages Sent Per Weekday"
            },
            legend: {
                display: false
            },
            maintainAspectRatio: false,
            scale: {
                ticks: {
                    suggestedMin: 0
                }
            }
        }
    });
}

CreateMessagesPerMonthChart();
CreateMessagesPerUserChart();
CreateMessagesPerWeekdayChart();
