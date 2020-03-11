



// 
// turn off all decoration initially
// 
$(document).ready(function() {
    var content = $(".leaf,.object,.array");
    content.css({
        'border-style': 'none',
        'background-color': 'transparent',
        'font-size': '18px'
    });
    content.show();
});

// 
// hide any broken image links 
// 
$("img").on("error", function() {
    $(this).hide();
});


// 
// handler for borders checkbox
// 
$('#cbx-borders').change(function() {
    var content = $(".leaf,.object,.array");
    if ($(this).is(':checked')) {
        content.css({ 'border-style': 'solid' });
    } else {
        content.css({ 'border-style': 'none' });
    }
    content.show();
});

// 
// handler for colours checkbox
// 
$('#cbx-colours').change(function() {
    var content = $(".leaf,.object,.array");
    if ($(this).is(':checked')) {
        
        var collections = $(".object,.array");
        collections.css({ 'background-color': 'rgba(252, 253, 175, 0.4)' });
        
        var leaves = $(".leaf");
        leaves.css({ 'background-color': 'rgba(208, 227, 204, 0.8)' });
        
        leaves.show();
        collections.show();

    } else {
        content.css({ 'background-color': 'transparent' });
    }
    content.show();
});






