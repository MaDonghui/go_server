$(document).ready(function() {
    $("#newPhone").on("submit", function(event){
        
        event.preventDefault();
        var data = {
            "image": `${document.getElementById("image").value}`,
            "brand": `${document.getElementById("brand").value}`,
            "model": `${document.getElementById("model").value}`,
            "os": `${document.getElementById("os").value}`,
            "screensize": parseInt(`${document.getElementById("screensize").value}`, 10)
        }
        
        $.ajax({
            url: "http://localhost:8000/new",
            type: "POST",
            data: JSON.stringify(data),
            success: function(result) {
                console.log(data)
                $("tr.items").remove();
                $.ajax({
                    type: "GET",
                    url: "http://localhost:8000/all",
                    success: function(result) {
                        //console.log(result);
                        $.each(result, function(i, product) {
                            $(".columnheaders").after(`
                                <tr class="items">
                                    <td data-label="Image">
                                        <img src="${product.image}" alt="${product.model} image">
                                    </td>
                                    <td data-label="Brand">${product.brand}</td>
                                    <td data-label="Model">${product.model}</td>
                                    <td data-label="Screen Size">${product.screensize}</td>
                                    <td data-label="Operating System">${product.os}</td>
                                </tr>`)
                        })
                    },
                    error: function(result) {
                        console.log("ERROR: " + result);
                    }
                });
            }
        })
    });
});
