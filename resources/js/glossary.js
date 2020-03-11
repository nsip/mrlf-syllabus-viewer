// 
// hide any broken image links 
// 
$("img").on("error", function() {
    $(this).hide();
});


// 
// very simple search for glossary items
//  - just looks for matching string.
// 
function searchFunction() {
    var input, filter, ul, li, a, i, x, txtValue;

    // get the search term
    input = document.getElementById("search");
    filter = input.value.toUpperCase();

    // get all the list items
    ul = document.getElementById("glossary-UL");
    li = ul.getElementsByTagName("li");
    for (i = 0; i < li.length; i++) {
        
        // check the content of all sub-spans of the li
        a = li[i].getElementsByTagName("span");

        var containsText = false;
        for (x = 0; x < a.length; x++) {
            txtValue = a[x].textContent || a[x].innerText;
            if (txtValue.toUpperCase().indexOf(filter) > -1) {
                containsText = true;
                break;
            }
        }

        // check for the string
        if (containsText) {
            li[i].style.display = "";
        } else {
            li[i].style.display = "none";
        }

    }
}