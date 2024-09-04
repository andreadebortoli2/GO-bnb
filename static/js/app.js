function sweetalertPrompt() {
    let toast = function (c) {
        const { message = "", icon = "success", position = "top-end" } = c;

        const Toast = Swal.mixin({
            toast: true,
            title: message,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.onmouseenter = Swal.stopTimer;
                toast.onmouseleave = Swal.resumeTimer;
            },
        });

        Toast.fire({});
    };

    let success = function (c) {
        const { message = "", title = "", footer = "" } = c;
        Swal.fire({
            icon: "success",
            title: title,
            text: message,
            footer: footer,
        });
    };

    let error = function (c) {
        const { message = "", title = "", footer = "" } = c;
        Swal.fire({
            icon: "error",
            title: title,
            text: message,
            footer: footer,
        });
    };

    async function custom(c) {
        const { icon = "", message = "", title = "", showConfirmButton = true } = c;

        const { value: result } = await Swal.fire({
            icon: icon,
            title: title,
            html: message,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById("start").value,
                    document.getElementById("end").value,
                ];
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            },
        });
        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== "") {
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
    };
}

function roomsLogics(roomID) {
    console.log(roomID);

    document.getElementById("check-availability").addEventListener("click", function () {
        let html = `
    <form class="row needs-validation" id="check-availability-form" action="reservation.html" method="post" novalidate>
        <div class="row py-3 mx-auto" id="reservation-dates-modal">
          <div class="col">
            <input class="form-control" type="text" name="start" id="start" placeholder="Arrival" disabled required />
          </div>
          <div class="col">
            <input class="form-control" type="text" name="end" id="end" placeholder="Departure" disabled required />
          </div>
        </div>
      </form>
    `;
        attention.custom({
            message: html,
            title: "Choose your dates",
            willOpen: () => {
                const elem = document.getElementById("reservation-dates-modal");
                const rangePicker = new DateRangePicker(elem, {
                    format: "yyyy-mm-dd",
                    showOnFocus: true,
                    minDate: new Date(),
                });
            },
            didOpen: () => {
                document.getElementById("start").removeAttribute("disabled");
                document.getElementById("end").removeAttribute("disabled");
            },
            callback: function (result) {
                console.log("called")

                let form = document.getElementById("check-availability-form");
                let formData = new FormData(form);
                formData.append("csrf_token", "{{.CSRFToken}}");
                formData.append("room_id", roomID);

                fetch('/search-availability-json', {
                    method: "post",
                    body: formData,
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.ok) {
                            attention.custom({
                                icon: 'success',
                                showConfirmButton: false,
                                msg: '<p>Room is available!</p>'
                                    + '<p><a href="/book-room?id='
                                    + data.room_id
                                    + '&s='
                                    + data.start_date
                                    + '&e='
                                    + data.end_date
                                    + '" class="btn btn-primary">Book now!</a></p>'
                            })
                        } else {
                            attention.error({
                                msg: "No availability",
                            })
                        }
                    })
            }
        });
    });
}