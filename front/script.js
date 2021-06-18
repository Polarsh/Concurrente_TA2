$('.form').find('input, textarea').on('keyup blur focus', function (e) {
  
  var $this = $(this),
      label = $this.prev('label');

	  if (e.type === 'keyup') {
			if ($this.val() === '') {
          label.removeClass('active highlight');
        } else {
          label.addClass('active highlight');
        }
    } else if (e.type === 'blur') {
    	if( $this.val() === '' ) {
    		label.removeClass('active highlight'); 
			} else {
		    label.removeClass('highlight');   
			}   
    } else if (e.type === 'focus') {
      
      if( $this.val() === '' ) {
    		label.removeClass('highlight'); 
			} 
      else if( $this.val() !== '' ) {
		    label.addClass('highlight');
			}
    }

});

$('.tab a').on('click', function (e) {
  
  e.preventDefault();
  
  $(this).parent().addClass('active');
  $(this).parent().siblings().removeClass('active');
  
  target = $(this).attr('href');

  $('.tab-content > div').not(target).hide();
  
  $(target).fadeIn(600);
  
});

///
`
const url = "http://localhost:9000"

var inputForm = document.getElementById("inputForm")

inputForm.addEventListener("submit", (e)=>{
  
  //prevent auto submission
  e.preventDefault()
  
  const formdata = new FormData(inputForm)
  fetch(url,{    
    method:"POST",
    body:formdata,
  }).then(
    response => response.text()
  ).then(
    (data) => {console.log(data);document.getElementById("gender").innerHTML=data}
  ).catch(
    error => console.error(error)
    )
}
      
)

`
    