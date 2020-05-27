document.addEventListener('DOMContentLoaded', () => {
    let url = window.location.pathname.split('/');
    let state = url[url.length-2];
    let city = url[url.length-1];

    fetch('/api/restaurantes/'+state+'/'+city)
    .then(res => {
        res.json().then(json => {
            if(json.length != 0){
                document.getElementById('title').innerText = 'Restaurantes em ' + json[0].address.city;
                document.getElementById('total_results').innerText = json.length + ' restaurantes encontrados';
        
                json.forEach(restaurant => {
                    restaurantCardInsert(restaurant)
                });
            } else {
                document.getElementById('title').innerText = 'Ops... Desculpe mas não foi encontrado nenhum restaurante.';    
            }
        });
    })
    .catch(err => {
        alert(err);
    })
})

document.getElementById("btn_filters").addEventListener("click", () => {   
    form = document.getElementById("form_filters")
    form_filter_name = document.getElementById("form_filter_name")

    if (form_filter_name.style.display !== "none") {
        form_filter_name.style.display = "none";
    }

    if (form.style.display === "none") {
        form.style.display = "block";
    } else {
        form.style.display = "none";
    }
})

document.getElementById("btn_filter_name").addEventListener("click", () => {   
    form = document.getElementById("form_filter_name")
    form_filters = document.getElementById("form_filters")
    document.getElementById('nome_busca').value = ''

    if (form_filters.style.display !== "none"){
        form_filters.style.display = "none";
    }

    if (form.style.display === "none") {
        form.style.display = "block";
    } else {
        form.style.display = "none";
    }
})

function restaurantCardInsert(restaurant) {
    restaurants_list = document.getElementById("restaurants_list");
    if(restaurant.average_price_per_meal == 56) {
        restaurant.average_price_per_meal = '55+'
    }
    let imagem = ``
    if(restaurant.plan.number != 0) {
        //imagem = `<img src="http://placehold.it/800x500?text=logo" width="800" height="400" class="card-img-top" alt="...">`
        url = ''
        restaurant.images.forEach(img => {
            if(img.main){
                url = img.url;
            }    
        });
        imagem = `<img src="${url}" width="800" height="400" class="card-img-top" alt="Restaurante logo">`
    }
    let horario = '';

    var groups = {};
    for (var i = 0; i < restaurant.opening_hours.length; i++) {
        var groupName = restaurant.opening_hours[i].opening_time + ' - ' + restaurant.opening_hours[i].closing_time;
        if (!groups[groupName]) {
            groups[groupName] = [];
        }
        groups[groupName].push(restaurant.opening_hours[i].weekday);
    }
    myArray = [];
    for (var groupName in groups) {
        myArray.push({hour: groupName, day: groups[groupName]});
    }

    myArray.forEach(element => {
        horario = horario + element.day.toString() + ': ' + element.hour + '. ';
    });

    let opcoes = `<div class="row options">`
    let count = 1;
    for(let key in restaurant.options) {
        if(restaurant.options.hasOwnProperty(key)) {
            let string = ''
            if(key == 'reservations' && restaurant.options[key] == true) {
                string = '<p class="card-text"><img src="/static/img/month.png" width="20px" height="20px"> Reservas</p>'
            } else if(key == 'parking' && restaurant.options[key] == true) {
                string = '<p class="card-text"><img src="/static/img/car.png" width="20px" height="20px"> Estacionamento</p>'
            } else if(key == 'takeout' && restaurant.options[key] == true) {
                string = '<p class="card-text"><img src="/static/img/box.png" width="20px" height="20px"> Embalagem para levar</p>'
            } else if(key == 'delivery' && restaurant.options[key] == true) {
                string = '<p class="card-text"><img src="/static/img/delivery.png" width="20px" height="20px"> Entregas</p>'
            } else if(key == 'wifi' && restaurant.options[key] == true) {
                string = '<p class="card-text"><img src="/static/img/wifi.png" width="20px" height="20px"> Wifi</p>'
            } else if(key == 'wheelchair' && restaurant.options[key] == true) {
                string = '<p class="card-text"><img src="/static/img/wheelchair.png" width="20px" height="20px"> Acesso para cadeirantes</p>'
            } else if(key == 'credit_card' && restaurant.options[key] == true) {
                string = '<p class="card-text"><img src="/static/img/credit.png" width="20px" height="20px"> Cartão de crédito</p>'
            } else if(key == 'outdoor' && restaurant.options[key] == true) {
                string = '<p class="card-text"><img src="/static/img/outdoor.png" width="20px" height="20px"> Ambiente exterior</p>'
            } else if(key == 'alcohol' && restaurant.options[key] == true) {
                string = '<p class="card-text"><img src="/static/img/beer.png" width="20px" height="20px"> Bebidas alcoolicas</p>'
            }
 
            if(string != '') {
                if(count <= 3) {
                    if(count == 1) {
                        opcoes = opcoes + '<div class="col-md-4">';
                    }
                    opcoes = opcoes + string;
                    count++;
                } else {
                    if(count == 4) {
                        opcoes = opcoes + '</div>';
                        opcoes = opcoes + '<div class="col-md-4">';
                        opcoes = opcoes + string;
                        count = 2;    
                    }
                }
            }
        }
    }
    opcoes = opcoes + '</div></div>'
    html = `
    <div class="card" style="margin-bottom: 15px;">
        ${imagem}
        
        <div class="card-body">
            <a href="/restaurantes/${restaurant.address.state_slug}/${restaurant.address.city_slug}/${restaurant.name_slug}"><h4 class="card-title">${restaurant.name}</h4></a>
            <div class="row">
                <div class="col-md-12">
                    <p class="card-text"><img src="/static/img/map-location.png" width="20px" height="20px"> ${restaurant.address.street}, ${restaurant.address.number} - ${restaurant.address.neighborhood}, ${restaurant.address.city} - ${restaurant.address.state}, ${restaurant.address.cep}</p>
                </div>
            </div>
            <div class="row">
                <div class="col-md-12 average-price">
                    <p class="card-text"><img src="/static/img/money.png" width="20px" height="20px"> Preço médio por refeição: R$${restaurant.average_price_per_meal}</p>
                </div>
            </div>
            <div class="row">
                <div class="col-md-12">
                    <p class="card-text"><img src="/static/img/clock.png" width="20px" height="20px"> Horário de funcionamento: ${horario}</p>
                </div>
            </div>
            <div class="row">
                <div class="col-md-12">
                    <p class="card-text">Opções</p>
                </div>
            </div>
            ${opcoes}
            <div class="row">
                <div class="col-md-12">
                    <hr>
                    <a href="/restaurantes/${restaurant.address.state_slug}/${restaurant.address.city_slug}/${restaurant.name_slug}" class="btn btn-primary">Ver mais informações</a>
                    <a href = "${restaurant.website}" class="btn btn-primary" target="_blank" rel=”noopener”>Visitar website</a>
                </div>
            </div>
        </div>
    </div>
    `
    restaurants_list.innerHTML = restaurants_list.innerHTML + html;
}


filter_form = document.getElementById("form_filters")

filter_form.addEventListener("submit", e => {
    e.preventDefault();

    const formData = new FormData(filter_form);
    let searchParams = new URLSearchParams();

    for(const pair of formData) {
        searchParams.append(pair[0], pair[1])
    }

    let url = window.location.pathname.split('/');
    let state = url[url.length-2];
    let city = url[url.length-1];

    fetch('/api/restaurantes/filtrar/'+state+'/'+city, {
        method: 'post',
        body: searchParams
    }).then(response => {
        response.json().then(json =>{
            if(json.length != 0){
                document.getElementById('title').innerText = 'Restaurantes em ' + json[0].address.city;
                document.getElementById('total_results').innerText = json.length + ' restaurantes encontrados';
                
                document.getElementById("restaurants_list").innerHTML = '';
                document.getElementById("btn_filters").click()
                json.forEach(restaurant => {
                    restaurantCardInsert(restaurant)
                });
            } else {
                document.getElementById("btn_filters").click()
                document.getElementById("restaurants_list").innerHTML = '<strong><p>Ops... Desculpe mas não foi encontrado nenhum restaurante.</p></strong>';   
            }
        })
    }).catch(err => {
        console.log(err)
    })

})

filter_name_form = document.getElementById("form_filter_name")

filter_name_form.addEventListener("submit", e => {
    e.preventDefault();

    const formData = new FormData(filter_name_form);
    let searchParams = new URLSearchParams();

    for(const pair of formData) {
        searchParams.append(pair[0], pair[1])
    }

    let url = window.location.pathname.split('/');
    let state = url[url.length-2];
    let city = url[url.length-1];

    fetch('/api/restaurantes/filtrar/nome/'+state+'/'+city, {
        method: 'post',
        body: searchParams
    }).then(response => {
        response.json().then(json =>{
            if(json.length != 0){
                document.getElementById('title').innerText = 'Restaurantes em ' + json[0].address.city;
                document.getElementById('total_results').innerText = json.length + ' restaurantes encontrados';
                
                document.getElementById("restaurants_list").innerHTML = '';
                document.getElementById("btn_filter_name").click()
                json.forEach(restaurant => {
                    restaurantCardInsert(restaurant)
                });
            } else {
                document.getElementById("btn_filter_name").click()
                document.getElementById("restaurants_list").innerHTML = '<strong><p>Ops... Desculpe mas não foi encontrado nenhum restaurante.</p></strong>';   
            }
        })
    }).catch(err => {
        console.log(err)
    })

})
