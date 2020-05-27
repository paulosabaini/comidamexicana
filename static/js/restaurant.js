document.addEventListener('DOMContentLoaded', () => {
    let url = window.location.pathname.split('/');
    let state = url[url.length-3];
    let city = url[url.length-2];
    let restaurant_slug = url[url.length-1];

    fetch(`/api/restaurantes/${state}/${city}/${restaurant_slug}`)
    .then(res => {
        res.json().then(json => {
            if(Object.keys(json).length!=0){       
                restaurantInsert(json)
                $('.carousel').carousel();
                createGoogleMap();
            } else {
                document.getElementById('title').innerText = 'Ops... Desculpe mas não foi encontrado nenhum restaurante.';    
            }
        });
    })
    .catch(err => {
        alert(err);
    })
})


function restaurantInsert(restaurant) {

    if(restaurant.average_price_per_meal == 56) {
        restaurant.average_price_per_meal = '55+'
    }

    whatsapp = ""
    if(restaurant.whatsapp) {
        whatsapp = `<a href="//api.whatsapp.com/send?phone=55${restaurant.mobilePhone.replace(' ', '').replace('(', '').replace(')', '').replace('-', '')}&text=Olá, encontrei seu restaurante pelo site www.comidamexicana.com.br" class="btn btn-outline-success" target="_blank" rel=”noopener”><img src="/static/img/whatsapp.png" width="20px" height="20px"> Enviar Mensagem</a>`
    }

    indicators = ""
    images = ""

    for(let i = 0; i < restaurant.images.length; i ++) {
        if(i <= restaurant.plan.number_of_images - 1) {
            active = ''
            if(i == 0){
                active = 'active'
            }
            indicators = indicators + `<li class="${active}" data-target="#carouselExampleIndicators" data-slide-to="${i}"></li>`
            images = images + `<div class="carousel-item ${active}"><img src="${restaurant.images[i].url}" class="d-block w-100" alt="imagem do restaurante"></div>`
        }
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

    video = ``

    if(restaurant.plan.number == 1 || restaurant.plan.number == 2) {
        if(restaurant.video_link != '') {
            let video_id = restaurant.video_link.split('v=')[1];
            let ampersandPosition = video_id.indexOf('&');
            if(ampersandPosition != -1) {
                video_id = video_id.substring(0, ampersandPosition);
            }
            video = `<div class="embed-responsive embed-responsive-16by9">
                        <iframe class="embed-responsive-item" src="https://www.youtube.com/embed/${video_id}" allowfullscreen></iframe>
                    </div>`
        }
    }

    social = ``

    if(restaurant.plan.number == 2){
        if(restaurant.instagram_link != '') {
            social += `<div class="col">
                            <a href="${restaurant.instagram_link}" target="_blank" rel=”noopener”><img src="/static/img/instagram.png" width="40px" height="40px"> </a>
                        </div>`
        }

        if(restaurant.facebook_link != '') {
            social += `<div class="col">
                            <a href="${restaurant.facebook_link}" target="_blank" rel=”noopener”><img src="/static/img/facebook.png" width="40px" height="40px"> </a>
                        </div>`
        }

        if(restaurant.twitter_link != '') {
            social += `<div class="col">
                            <a href="${restaurant.twitter_link}" target="_blank" rel=”noopener”><img src="/static/img/twitter.png" width="40px" height="40px"> </a>
                        </div>`
        }
    }

    rating_showcase = ``

    if(restaurant.plan.number == 2){
        if(restaurant.opinion_phrase != ''){
            rating_showcase = `<div class="alert alert-info showcase" role="alert">
                                    <p>${restaurant.opinion_phrase}</p>
                                    <p>- ${restaurant.opinion_author}</p>
                                </div>`
        }
    }

    ifood = ``

    if(restaurant.plan.number == 2){
        if(restaurant.ifood_link != '') {
            ifood = `<a href="${restaurant.ifood_link}" target="_blank" rel=”noopener”>Peça no <img src="/static/img/ifood.png" width="60px" height="60px"></a>`  
        }    
    }

    html = `
    <h1 id="restaurant-name" class="restaurant-name">${restaurant.name}</h1>
    <p class="card-text"><img src="/static/img/map-location.png" width="20px" height="20px"> <span id="address">${restaurant.address.street}, ${restaurant.address.number} - ${restaurant.address.neighborhood}, ${restaurant.address.city} - ${restaurant.address.state}, ${restaurant.address.cep}</span></p>
    <div id="map-location">
        <div id="map" style="height: 300px; width: 100%;"></div>
        <a class="btn btn-primary ml-auto direction" id="direction" target="_blank" rel="noopener">Como Chegar</a>
    </div>
    <div id="contact" class="contact">
        <div class="row">
            <div class="col">
                ${restaurant.phone ? `<p id="phone"><img src="/static/img/phone.png" width="20px" height="20px"> ${restaurant.phone}</p>` : ''}
            </div>
            <div class="col">
                ${restaurant.mobilePhone ? `<p id="mobilePhone"><img src="/static/img/cellphone.png" width="20px" height="20px"> ${restaurant.mobilePhone}</p>` : ''}
            </div>
        </div>
        <div class="row">
            <div class="col">
                <a class="btn btn-outline-info" href="${restaurant.website}" target="_blank" rel=”noopener”>Visitar Website</a>
            </div>
            <div class="col">
                ${whatsapp}    
            </div>
        </div>
        <div class="row">
            <div class="col">
                ${ifood}
            </div>    
        </div>
    </div>
    <div id="carouselExampleIndicators" class="carousel slide" data-ride="carousel">
        <ol class="carousel-indicators">
        ${indicators}
        </ol>
        <div class="carousel-inner">
        ${images}
        </div>
    </div>

    <div id="opening-time" class="opening-time">
        <p class="card-text"><img src="/static/img/clock.png" width="20px" height="20px"> <strong>Horário de funcionamento:</strong> ${horario}</p>
    </div>
    <div id="price" class="price">
        <p class="card-text"><img src="/static/img/money.png" width="20px" height="20px"> <strong>Preço médio por refeição:</strong> R$${restaurant.average_price_per_meal}</p>
    </div>
    <div id="options" class="options">
        <p><strong>Opções:</strong></p>
        ${opcoes}
    </div>
    <hr>
    <div id="description">
        <p id="description-text" class="description-text">${restaurant.description.replace(/(?:\r\n|\r|\n)/g,'<br>')}</p>
    </div>
    <hr>
    ${video}
    <div id="social" class="social">
        <div class="row">
            ${social}
        </div>    
    </div>
    <div id="rating-showcase" class="rating-showcase">
        ${rating_showcase}
    </div>
    <div id="promotion-image">
        <img src="/static/img/promo.png" width="100%" height="170">
    </div>
    <div>
        
    </div>
    `

    document.getElementById('restaurant-content').innerHTML = html
}

function createGoogleMap() {
    let map = new google.maps.Map(document.getElementById('map'), {
        zoom: 15,
        center: {lat: -34.397, lng: 150.644}
    });

    let address = document.getElementById('address').innerText;

    let geocoder = new google.maps.Geocoder();

    let direction = 'https://www.google.com/maps/dir/?api=1&'

    geocoder.geocode({'address': address}, function(results, status) {
        if (status === 'OK') {
          map.setCenter(results[0].geometry.location);
          var marker = new google.maps.Marker({
            map: map,
            position: results[0].geometry.location
          });
          direction += 'destination='+results[0].geometry.location.toString().replace('(','').replace(')','').replace(' ','')
          document.getElementById("direction").href = direction;
        } else {
          alert('Geocode was not successful for the following reason: ' + status);
        }
    });
}