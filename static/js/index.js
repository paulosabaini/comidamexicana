document.addEventListener('DOMContentLoaded', () => {
    fetch("/api/cities").then((response)=>{
        response.json().then(json =>{
            let list_results = document.getElementById('list_results')       
            document.getElementById('cidade').addEventListener('keyup', () => {
                input = document.getElementById('cidade');

                let min_characters = 3;

                if (input.value.length < min_characters){
                    list_results.innerHTML = '';
                    return
                } else {
                    list_results.innerHTML = '';
                    let firstCity = true;
                    let cityNotFound = true; 
                    json.forEach((item) => {
                        //const lowerValue = input.value.toLowerCase()
                        //const lowerCity = item.city.toLowerCase()
                        const lowerValue = removeAcento(input.value)
                        const lowerCity = removeAcento(item.city)
                        
                        if (lowerCity.includes(lowerValue)) {
                            cityNotFound = false;
                            let div = document.createElement('div');
                            if (firstCity){
                                div.className = 'autocompleteSuggestion autocompleteSelected';    
                            } else {
                                div.className = 'autocompleteSuggestion';
                            }
                            let city_name = document.createElement('p');
                            let city_slug = document.createElement('span');
                            let state_slug = document.createElement('span');
                            let icon = document.createElement('img');

                            icon.className = 'locationIcon';
                            icon.src = '/static/img/location.png';
                            city_name.textContent = item.city + ' - ' + item.state;
                            city_slug.style = 'display: none;';
                            state_slug.style = 'display: none;';
                            city_slug.className = 'citySlug';
                            state_slug.className = 'stateSlug';
                            city_slug.textContent = item.city_slug;
                            state_slug.textContent = item.state_slug;

                            div.appendChild(icon);
                            div.appendChild(city_name);
                            div.appendChild(city_slug);
                            div.appendChild(state_slug);

                            div.onclick = (e) => {
                                if(e.target.parentNode.innerText != '') {
                                    input.value = e.target.parentNode.innerText.replace(' -',',');
                                    document.getElementById("list_results").innerHTML = '';
                                    restaurantsSearch(e.target.parentNode);
                                }
                            }
                            list_results.appendChild(div);
                            firstCity = false;
                        }
                    })

                    if (cityNotFound) {
                        let div = document.createElement('div');
                        div.className = 'autocompleteSuggestion';
                        let msg = document.createElement('p');
                        msg.textContent = 'Ops... essa cidade não possui nenhum restaurante cadastrado. Nos envie uma mensagem informando qual a sua cidade que vamos em busca de restaurantes para você.'
                        div.appendChild(msg);
                        list_results.appendChild(div);    
                    }
                }
            })  
        })
    })

    function removeAcento (text)
    {       
        text = text.toLowerCase();                                                         
        text = text.replace(new RegExp('[ÁÀÂÃ]','gi'), 'a');
        text = text.replace(new RegExp('[ÉÈÊ]','gi'), 'e');
        text = text.replace(new RegExp('[ÍÌÎ]','gi'), 'i');
        text = text.replace(new RegExp('[ÓÒÔÕ]','gi'), 'o');
        text = text.replace(new RegExp('[ÚÙÛ]','gi'), 'u');
        text = text.replace(new RegExp('[Ç]','gi'), 'c');
        return text;                 
    }

    fetch('/api/restaurantes/2').then(response => {
        response.json().then(json =>{
            if(json.length != 0){
                json.forEach(restaurant => {
                    restaurantCardInsert(restaurant)
                });
            } else {
                restaurantCardInsert({name: 'Restaurante em destaque',name_slug: '', id: '', website: 'https://www.comidamexicana.com.br/listar-negocio', images: [{url: '/static/img/placeholder.png', main: 'true'}], address: {state_slug: '', city_slug: ''}})
                restaurantCardInsert({name: 'Restaurante em destaque',name_slug: '', id: '', website: 'https://www.comidamexicana.com.br/listar-negocio', images: [{url: '/static/img/placeholder.png', main: 'true'}], address: {state_slug: '', city_slug: ''}})
                restaurantCardInsert({name: 'Restaurante em destaque',name_slug: '', id: '', website: 'https://www.comidamexicana.com.br/listar-negocio', images: [{url: '/static/img/placeholder.png', main: 'true'}], address: {state_slug: '', city_slug: ''}})
            }
        })
    }).catch(err => {
        console.log(err)    
    })

})

document.addEventListener("click", (evt) => {
    const flyoutElement = document.getElementById("suggestions");
    let targetElement = evt.target; // clicked element

    do {
    if (targetElement == flyoutElement) {
        return;
    }
    // Go up the DOM
    targetElement = targetElement.parentNode;
    } while (targetElement);

    // This is a click outside.
    if (document.getElementById("list_results").innerHTML != '' && evt.target.id != 'btn-submit'){
        document.getElementById("list_results").innerHTML = ''
    }
});


search_form = document.getElementById('search-form');
search_form.addEventListener('submit', (event) => {
    event.preventDefault();
    let suggestions = document.getElementById('suggestions');
    element = suggestions.querySelector('.autocompleteSelected');
    if (element) {
        input = document.getElementById('cidade');
        input.value = element.innerText.replace(' -',',');
        restaurantsSearch(element);
        document.getElementById("list_results").innerHTML = '';
    }
})

function restaurantsSearch(search) {
    let state_slug = search.querySelector('.stateSlug').innerText;
    let city_slug = search.querySelector('.citySlug').innerText;
    window.location.href = '/restaurantes/'+state_slug+'/'+city_slug;
}

newsletter_form = document.getElementById('form-newsletter')
newsletter_form.addEventListener('submit', e => {
    e.preventDefault();
    newsl = {
        email: newsletter_form.email.value
    }
    fetch('/newsletter/create', {
        method: 'post',
        body: JSON.stringify(newsl)
    }).then(response => {
        document.getElementById('email').value = '';
        alert("Email enviado com sucesso!")
    }).catch(err => {
        console.log(err)
    })    
})

function restaurantCardInsert(restaurant) {
    destaques = document.getElementById('destaques')

    url = ''
        restaurant.images.forEach(img => {
            if(img.main){
                url = img.url;
            }    
        });

    html = `
            <div class="col-md-4 rest-dest">
                <div class="card" style="width: 18rem;">
                    <img class="card-img-top" src="${url}" alt="restaurante logo" width="60" height="250">
                    <div class="card-body">
                    <h5 class="card-title">${restaurant.name}</h5>
                    <a href="/restaurantes/${restaurant.address.state_slug}/${restaurant.address.city_slug}/${restaurant.name_slug}" class="btn btn-success">Ver mais</a>
                    <a href="${restaurant.website}" class="btn btn-success" target="_blank" rel=”noopener”>Website</a>
                    </div>
                </div>
            </div>`

    destaques.innerHTML += html

}