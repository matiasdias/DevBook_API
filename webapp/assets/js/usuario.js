$('#parar-de-seguir').on('click', pararDeSequir);
$('#seguir').on('click', seguir);
$('#editar-usuario').on('submit', editar);
$('#atualizar-senha').on('submit', atualizarSenha);
$('#deletar-usuario').on('click', deletarUsuario);


//pararDeSequir um usuario
function pararDeSequir() {
    const usuarioId = $(this).data('usuario-id');
    $(this).prop('disabled', true)

    $.ajax({
        url: `/usuarios/${usuarioId}/parar-de-seguir`,
        method: "POST"
    }).done(function() {
        window.location = `/usuarios/${usuarioId}`;
    }).fail(function() {
        Swal.fire("Ops...", "Erro ao parar de seguir o usuario!", "error");
        $('#parar-de-seguir').prop('disabled', false);
    });

}
//Seguir um usuario novo 
function seguir() {
    const usuarioId = $(this).data('usuario-id');
    $(this).prop('disabled', true)

    $.ajax({
        url: `/usuarios/${usuarioId}/seguir`,
        method: "POST"
    }).done(function() {
        window.location = `/usuarios/${usuarioId}`;
    }).fail(function() {
        Swal.fire("Ops...", "Erro ao seguir o usuario!", "error");
        $('#seguir').prop('disabled', false);
    });

}
//editar edita o usuario logado 
function editar(evento) {
    evento.preventDefault();

    $.ajax({
        url: "/editar-usuario",
        method: "PUT",
        data: {
            nome: $('#nome').val(),
            email: $('#email').val(),
            nick: $('#nick').val(),
        }
    }).done(function() {
        Swal.fire("Sucesso...", "Usuario editado!", "success")
            .then(function() {
                window.location = "/perfil";
            });
    }).fail(function () {
        Swal.fire("OPS...", "ERRO AO EDITAR O USUARIO !", "error")
    })
}

function atualizarSenha(evento) {
    evento.preventDefault();

    if ($('#nova-senha').val() != $('#confirmar-senha').val()) {
        Swal.fire("Ops...", "As senhas não coemcidem", "warning");
        return  
    }

    $.ajax({
        url: "/atualizar-senha",
        method: "POST",
        data: {
            atual: $('#senha-atual').val(),
            nova: $('#nova-senha').val()
        }
    }).done(function() {
        Swal.fire("Sucesso", "A senha foi atualizada", "success")
            .then(function() {
                    window.location = "/perfil";
            })
    }).fail(function() {
        Swal.fire("Ops...", "Erro ao atualizar a senha", "error");
    });

}

function deletarUsuario() {
    Swal.fire({
        title: "Atenção",
        text: "Tem certeza que deseja apagar a sua conta? Essa é um ação irreversivel",
        showCancelButton: true,
        cancelButtonText: "Cancelar",
        icon: "warning"
    }).then(function(confirmacao) {
        if (confirmacao.value) {
            $.ajax({
                url: "/deletar-usuario",
                method: "DELETE"
            }).done(function() {
                Swal.fire("Sucesso", "Usuario excluido com sucesso", "success")
                    .then(function() {
                        window.location = "/logout";
                    })
            }).fail(function() {
                Swal.fire("Ops...", "Erro ao excluir o usuario", "error");
            });
        }
    })
}

