var URL_ROOT = 'http://127.0.0.1:8000';

var vueutil = {
    post: function (vuectrl, postUrl, postData, handler) {
        this.laddaStart();
        var postData = vuectrl.srvModel;
        vuectrl.modelstate = {};
        axios({
            method: 'POST',
            url: URL_ROOT + postUrl,
            data: new URLSearchParams(postData).toString(),
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }).then(response => {
            if (response.status === 201) {
                window.location = response.headers.location;
            }
            handler(response);
        }).catch(error => {
            if (error.response.status === 400) {
                vuectrl.modelstate = error.response.data;
            } else {
                if (error.response.status >= 500 && error.response.status <= 599) {
                    alert('Server Error!');
                }

                console.log(error);
            }
        }).finally(() => {
            this.laddaStop()
        });
    },

    laddaStart: function () {
        var button = document.querySelector('button[type=submit]');
        if (!button.hasAttribute('data-style')) {
            button.setAttribute('data-style', 'zoom-out');
        }
        var l = Ladda.create(button);
        l.start();
    },

    laddaStop: function () {
        Ladda.stopAll();
    }
};