states = [
    { uf: 'AC', nome: 'Acre' },
    { uf: 'AL', nome: 'Alagoas' },
    { uf: 'AP', nome: 'Amapá' },
    { uf: 'AM', nome: 'Amazonas' },
    { uf: 'BA', nome: 'Bahia' },
    { uf: 'CE', nome: 'Ceará' },
    { uf: 'DF', nome: 'Distrito Federal' },
    { uf: 'ES', nome: 'Espirito Santo' },
    { uf: 'GO', nome: 'Goiás' },
    { uf: 'MA', nome: 'Maranhão' },
    { uf: 'MS', nome: 'Mato Grosso do Sul' },
    { uf: 'MT', nome: 'Mato Grosso' },
    { uf: 'MG', nome: 'Minas Gerais' },
    { uf: 'PA', nome: 'Pará' },
    { uf: 'PB', nome: 'Paraíba' },
    { uf: 'PR', nome: 'Paraná' },
    { uf: 'PE', nome: 'Pernambuco' },
    { uf: 'PI', nome: 'Piauí' },
    { uf: 'RJ', nome: 'Rio de Janeiro' },
    { uf: 'RN', nome: 'Rio Grande do Norte' },
    { uf: 'RS', nome: 'Rio Grande do Sul' },
    { uf: 'RO', nome: 'Rondônia' },
    { uf: 'RR', nome: 'Roraima' },
    { uf: 'SC', nome: 'Santa Catarina' },
    { uf: 'SP', nome: 'São Paulo' },
    { uf: 'SE', nome: 'Sergipe' },
    { uf: 'TO', nome: 'Tocantins' }
]

function populate(selector) {
    var select = $(selector);
    var hours, minutes, ampm;
    for (var i = 0; i <= 1420; i += 30) {
        hours = Math.floor(i / 60);
        minutes = i % 60;
        if (minutes < 10) {
            minutes = '0' + minutes; // adding leading zero to minutes portion
        }
        //add the value to dropdownlist
        if(i == 0){
            select.append($('<option></option>')
            .attr('value', -1)
            .text(''));
        }
        select.append($('<option></option>')
            .attr('value', hours + ':' + minutes)
            .text(hours + ':' + minutes));
    }
}

//Calling the function on pageload
window.onload = function (e) {
    elements = document.getElementsByClassName("timeDropdownlist");
    for(i=0;i<elements.length;i++){
        populate('#'+elements[i].id);
    }

    let path = window.location.href.split('/');
    let pathName = path[path.length - 1]

    if(pathName.includes("gratis")) {
        this.document.getElementById('chbPlanGratis').checked = true;
        document.getElementById('youtube').readOnly = true;
        document.getElementById('instagram').readOnly = true; 
        document.getElementById('facebook').readOnly = true; 
        document.getElementById('twitter').readOnly = true;
    } else if(pathName.includes("intermediario")) {
        this.document.getElementById('chbPlanInterm').checked = true;
        document.getElementById('youtube').readOnly = false;
        document.getElementById('instagram').readOnly = true; 
        document.getElementById('facebook').readOnly = true; 
        document.getElementById('twitter').readOnly = true; 
    } else if(pathName.includes("completo")) {
        this.document.getElementById('chbPlanComp').checked = true;
        document.getElementById('youtube').readOnly = false;
        document.getElementById('instagram').readOnly = false; 
        document.getElementById('facebook').readOnly = false; 
        document.getElementById('twitter').readOnly = false; 
    }

    id('telefoneCel').onkeyup = function(){
        mascara( this, mtel );
    }
    id('telefoneFixo').onkeyup = function(){
        mascara( this, mtel );
    }
}

function limpa_formulario_cep() {
    //Limpa valores do formulário de cep.
    document.getElementById('rua').value=("");
    document.getElementById('bairro').value=("");
    document.getElementById('cidade').value=("");
    document.getElementById('uf').value=("");
}

function meu_callback(conteudo) {
if (!("erro" in conteudo)) {
    //Atualiza os campos com os valores.
    document.getElementById('rua').value=(conteudo.logradouro);
    document.getElementById('bairro').value=(conteudo.bairro);
    document.getElementById('cidade').value=(conteudo.localidade);
    for (let i=0; i < states.length; i++)
        if (states[i]['uf'] == conteudo.uf)
            document.getElementById('uf').value = states[i]['nome']
    //document.getElementById('uf').value=(conteudo.uf);
    document.getElementById('numero').focus();
} //end if.
else {
    //CEP não Encontrado.
    limpa_formulario_cep();
    alert("CEP não encontrado.");
}
}

function pesquisacep(valor) {

//Nova variável "cep" somente com dígitos.
var cep = valor.replace(/\D/g, '');

//Verifica se campo cep possui valor informado.
if (cep != "") {

    //Expressão regular para validar o CEP.
    var validacep = /^[0-9]{8}$/;

    //Valida o formato do CEP.
    if(validacep.test(cep)) {

        //Preenche os campos com "..." enquanto consulta webservice.
        document.getElementById('rua').value="...";
        document.getElementById('bairro').value="...";
        document.getElementById('cidade').value="...";
        document.getElementById('uf').value="...";

        //Cria um elemento javascript.
        var script = document.createElement('script');

        //Sincroniza com o callback.
        script.src = 'https://viacep.com.br/ws/'+ cep + '/json/?callback=meu_callback';

        //Insere script no documento e carrega o conteúdo.
        document.body.appendChild(script);

    } //end if.
    else {
        //cep é inválido.
        limpa_formulario_cep();
        alert("Formato de CEP inválido.");
    }
} //end if.
else {
    //cep sem valor, limpa formulário.
    limpa_formulario_cep();
}
};

document.getElementById('chbPlanGratis').addEventListener('change', () => {
    if(document.getElementById('chbPlanGratis').checked){
        document.getElementById('youtube').readOnly = true;
        document.getElementById('instagram').readOnly = true; 
        document.getElementById('facebook').readOnly = true; 
        document.getElementById('twitter').readOnly = true;
        document.getElementById('ifood').readOnly = true;        
    }
})

document.getElementById('chbPlanInterm').addEventListener('change', () => {
    if(document.getElementById('chbPlanInterm').checked){
        document.getElementById('youtube').readOnly = false;
        document.getElementById('instagram').readOnly = true; 
        document.getElementById('facebook').readOnly = true; 
        document.getElementById('twitter').readOnly = true;
        document.getElementById('ifood').readOnly = true;         
    }
})

document.getElementById('chbPlanComp').addEventListener('change', () => {
    if(document.getElementById('chbPlanComp').checked){
        document.getElementById('youtube').readOnly = false;
        document.getElementById('instagram').readOnly = false; 
        document.getElementById('facebook').readOnly = false; 
        document.getElementById('twitter').readOnly = false;
        document.getElementById('ifood').readOnly = false;         
    }
})

/* Máscaras ER */
function mascara(o,f){
    v_obj=o
    v_fun=f
    setTimeout("execmascara()",1)
}

function execmascara(){
    v_obj.value=v_fun(v_obj.value)
}

function mtel(v){
    v=v.replace(/\D/g,""); //Remove tudo o que não é dígito
    v=v.replace(/^(\d{2})(\d)/g,"($1) $2"); //Coloca parênteses em volta dos dois primeiros dígitos
    v=v.replace(/(\d)(\d{4})$/,"$1-$2"); //Coloca hífen entre o quarto e o quinto dígitos
    return v;
}

function id( el ){
    return document.getElementById( el );
}

document.getElementById('cep').addEventListener('keyup', () => {
    value = document.getElementById('cep').value

    if (value.length == 5) {
        value = value + '-'
        document.getElementById('cep').value = value;
    } 
})

function validateForm() {

    if(document.formEnvio.telefoneFixo.value == '' && document.formEnvio.telefoneCel.value == '') {
        alert("Informe ao menos um telefone para contato.")
        return false;
    }

    if(!document.formEnvio.seg.checked && !document.formEnvio.ter.checked && !document.formEnvio.qua.checked && 
        !document.formEnvio.qui.checked && !document.formEnvio.sex.checked && !document.formEnvio.sab.checked && 
        !document.formEnvio.dom.checked) {
        alert("Selecione o horário de funcionamento.")
        return false;
    }

    if(!document.formEnvio.reserva.checked && !document.formEnvio.estacionamento.checked && !document.formEnvio.embalagem.checked && 
        !document.formEnvio.entrega.checked && !document.formEnvio.wifi.checked && !document.formEnvio.cadeirantes.checked && 
        !document.formEnvio.cartaocredito.checked && !document.formEnvio.exterior.checked && !document.formEnvio.bebidas.checked) {
            alert("Selecione alguma das opções disponíveis.")
            return false;
    }

    return true

}