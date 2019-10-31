// Note about canvas sizes:
// https://www.chartjs.org/docs/latest/general/responsive.html#important-note

let rgbColour = "rgb(0,198,255)";

async function createMessagesPerMonthChart(jsonData) {
    let data = {
        labels: [],
        datasets: [{
            data: [],
            borderColor: rgbColour,
            lineTension: 0.25
        }]
    };

    for (const year in jsonData.msgsPerMonth) {
        for (const month in jsonData.msgsPerMonth[year]) {
            data.labels.push(`${year}-${month}`);
            data.datasets[0].data.push(jsonData.msgsPerMonth[year][month]);
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

async function createMessagesPerUserChart(jsonData) {
    let data = {
        labels: [],
        datasets: [{
            data: [],
            backgroundColor: rgbColour
        }]
    };

    for (const user in jsonData.msgsPerUser) {
        data.labels.push(user);
        data.datasets[0].data.push(jsonData.msgsPerUser[user]);
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

async function createMessagesPerWeekdayChart(jsonData) {
    let data = {
        labels: [],
        datasets: [{
            data: [],
            borderColor: rgbColour
        }]
    };

    let weekdays = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"];
    weekdays.forEach((weekday) => {
        data.labels.push(weekday);
        data.datasets[0].data.push(jsonData.msgsPerWeekday[weekday]);
    });

    const ctx = document.getElementById("messagesPerWeekdayChart").getContext("2d");
    new Chart(ctx, {
        type: "radar",
        data: data,
        options: {
            title: {
                display: true,
                text: "Total Messages Sent Per Weekday"
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

async function setTitle(jsonData) {
    let title = document.querySelector("#title");
    title.textContent = `Messenger Stats for conversation: ${jsonData.title}`
}

window.addEventListener("load", async () => {

    let currentUrl = new URL(window.location.href);
    let id = currentUrl.searchParams.get("id");

    if (id !== null) {
        let rawData = await fetch(`/api/data?id=${id}`);
        let data = await rawData.json();

        setTitle(data);
        createMessagesPerMonthChart(data);
        createMessagesPerUserChart(data);
        createMessagesPerWeekdayChart(data);
    }

});
