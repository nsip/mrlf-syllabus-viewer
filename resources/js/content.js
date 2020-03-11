



// 
// some manipulations required after json
// has been processed
// 
$(document).ready(function() {
    // 
    // remove spurious empty nodes
    // 
    var content = $(".leaf, span:contains([])");
    content.remove();
    // 
    // insert section titles
    // 
    var objectives = $(".objectives");
    objectives.before("<h2>Objectives</h2>");

    // var stages = $(".container > .children");
    // stages.before("<h2>Stages</h2>");

    var concepts = $(".container > .children > .children");
    concepts.before("<h2>Concepts</h2>");
});

// 
// hide any broken image links 
// 
$("img").on("error", function() {
    $(this).hide();
});

