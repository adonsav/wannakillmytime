// Prompt is a Javascript module for all alerts, notifications and custom popups
function Prompt() {
    let toast = function(params) {
        const {
            title = '',
            icon = 'success',
            position = 'top-end',
        } = params;

        const Toast = Swal.mixin({
            toast: true,
            title: title,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let success = function(params) {
        const {
            text = '',
            title = '',
            footer = '',
        } = params

        Swal.fire({
            icon: 'success',
            title: title,
            text: text,
            footer: footer,
        })
    }

    let error = function(params) {
        const {
            text = '',
            title = '',
            footer = '',
        } = params

        Swal.fire({
            icon: 'error',
            title: title,
            text: text,
            footer: footer,
            showDenyButton: true,
            denyButtonText: "Close",
            showConfirmButton: false,
        })
    }

    async function custom(c) {
        const {
            html = '',
            title = '',
        } = c;

        const { value: result } = await Swal.fire({
            title: title,
            html: html,
            width: '22em',
            backdrop: false,
            focusConfirm: false,
            confirmButtonColor: '#0D6EFD',
            showCancelButton: true,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            didOpen: () => {
                if (c.didOpen !== undefined){
                    c.didOpen();
                }
            }
        })

        // result is whatever the user clicked on(it is passed to us by SweetAlert)
        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== '') {
                    if (c.callback !== undefined) {
                        c.callback(result);
                    }
                } else {
                    c.callback(false);
                }
            } else {
                c.callback(false);
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}