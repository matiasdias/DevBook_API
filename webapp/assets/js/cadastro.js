$('#formulario-cadastro').on('submit', criarUsuario);

function criarUsuario(evento) {
    evento.preventDefault();

    if ($('#senha').val() != $('#confirmar-senha').val()) {
        Swal.fire("Ops......", "As senhas são diferentes", "error");
        return;
    }

    $.ajax({
        url: "/usuarios",
        method: "POST",
        data: {
            nome: $('#nome').val(),
            email: $('#email').val(),
            nick: $('#nick').val(),
            senha: $('#senha').val()
        }
    }).done(function() {
        Swal.fire("Sucesso!", "Usuario cadastrado com sucesso", "success")
        .then(function() {
            $.ajax({
                url: "/login",
                method: "POST",
                data: {
                    email: $('#email').val(),
                    senha: $('#senha').val()
                }
            }).done(function() {
                window.location = "/home";
            }).fail(function() {
                Swal.fire("Ops......", "Erro ao autenticar o usuario", "error");
            })
        })
    }).fail(function() {
        Swal.fire("Ops......", "Erro ao cadastrar o usuario", "error");
    });
}