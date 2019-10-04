function openCity(evt, cityName) {
    var i, tabcontent, tablinks;
    tabcontent = document.getElementsByClassName("tabcontent");
    for (i = 0; i < tabcontent.length; i++) {
      tabcontent[i].style.display = "none";
    }
    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
      tablinks[i].className = tablinks[i].className.replace(" active", "");
    }
    document.getElementById(cityName).style.display = "block";
    evt.currentTarget.className += " active";
  }
  
  var acc = document.getElementsByClassName("accordion");
  var i;
  
  for (i = 0; i < acc.length; i++) {
    acc[i].addEventListener("click", function() {
      this.classList.toggle("active");
      var panel = this.nextElementSibling;
      if (panel.style.display === "block") {
        panel.style.display = "none";
      } else {
        panel.style.display = "block";
      }
    });
  }



  $("#submitBtn").click(function(event){
    event.preventDefault();

    var gname = $("#groupname").val();

    var checked = [];
    $.each($('input[type="checkbox"]:checked'), function(){            
        checked.push($(this).val());
    });

   var numChecked = $('input[type="checkbox"]:checked').length

     var format = checked.join(',');

    $.post("/contact/group", 

    {
      gname: gname,
      format: format,
      numChecked:numChecked,
        
    }, function(data, status){
      window.location.href="/contact/group/"+data
    })
  })

  // Get the modal
var modal = document.getElementById('msgModal');
window.onclick = function(event) {
    if (event.target == modal) {
        modal.style.display = "none";
    }
}