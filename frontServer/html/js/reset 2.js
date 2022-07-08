$(document).ready(function() {
    $("#reset").click(function() {
        $.ajax({
            type: "DELETE",
            url: "http://localhost:8000/reset",
            success: function(result) {
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
            },
            error: function(result) {
                console.log("ERROR: " + result);
            }
        });
    });
});