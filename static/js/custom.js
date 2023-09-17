$.MyAPP = {};

$.MyAPP.Disconnect = function (cname){
  console.log(cname)
  $.ajax({
    type: "DELETE",
    dataType: "json",
    url: "api/v1/session",
    data: JSON.stringify({ "cname": cname }),
    success: function(data) {
      location.reload();
      console.log(data);
    },
    error: function(a,b,c) {
      console.log(a,b,c)
      location.reload();
    }
  });
}

$.MyAPP.Restart = function (sname){
  console.log(sname)
  $.ajax({
    type: "DELETE",
    dataType: "json",
    url: "/signal",
    data: JSON.stringify({ "sname": sname }),
    success: function(data) {
      location.reload();
      console.log(data);
    },
    error: function(a,b,c) {
      console.log(a,b,c)
      location.reload();
    }
  });
}

// Copy modal content to clipboard
$('.copy-modal-data-btn').on('click', function () {
var $modal = $(this).closest('.modal');
var content = '';
$modal.find('.copy-details').each(function () {
content += $(this).text().replace(/:\s/g, ': ') + '\n';
});
content = content.replace(/:\s\n/g, ': ');
var $tempTextArea = $('<textarea>').val(content).css('position', 'absolute').css('left', '-9999px');
$('body').append($tempTextArea);
$tempTextArea.select();
document.execCommand('copy');
$tempTextArea.remove();
});

$(function() {
  new Clipboard('.button-copy');

  //$( ".btn-disconnect" ).click(function() {
  //  alert( "Handler for .click() called." );
  //});
  //window.location.reload();
})

function createEditor(name, size, theme, mode, readonly) {
  // find the textarea
  var textarea = document.querySelector("form textarea[name=" + name + "]");

  // create ace editor 
  var editor = ace.edit()
  editor.container.style.height = size

  editor.setTheme("ace/theme/" + theme); //"clouds_midnight"
  //editor.setTheme("ace/theme/twilight");
  //editor.setTheme("ace/theme/iplastic");

  editor.session.setMode("ace/mode/" + mode);

  editor.setReadOnly(readonly);
  editor.setShowPrintMargin(false);
  editor.session.setUseWrapMode(true);
  editor.session.setValue(textarea.value)
  // replace textarea with ace
  textarea.parentNode.insertBefore(editor.container, textarea)
  textarea.style.display = "none"
  // find the parent form and add submit event listener
  var form = textarea
  while (form && form.localName != "form") form = form.parentNode
  form.addEventListener("submit", function() {
      // update value of textarea to match value in ace
      textarea.value = editor.getValue()
  }, true)
}