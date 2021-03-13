const Controller = {
    search: (ev) => {
        ev.preventDefault();
        const form = document.getElementById("form");
        const data = Object.fromEntries(new FormData(form));
        const response = fetch(`/search?q=${data.query}`).then((response) => {
            response.json().then((results) => {
                Controller.updateTable(results);
            });
        });
    },

    updateTable: (results) => {
        const table = document.getElementById("search-results");
        const rows = [];
        for (let result of results) {
            rows.push(`<li class="list-group-item">${result}</li>`);
        }
        table.innerHTML = rows.join("");
    },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
