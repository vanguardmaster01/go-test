$(document).ready(function() {
    $('#register').on('click', function(){
        $.ajax({
            type:"POST",
            url : '/products',
            data:{
                name : $('#register_name').val(),
                description : $('#register_description').val(),
                price : $('#register_price').val(),
            },
            dataType : 'JSON',
            success : function(response){
                if(response.status == 'success'){
                    $('#registerModal').modal('hide');
                }
               console.log("success");
            },
            error: function(response){
               
                console.log(response);
            }
        });
        
        // $("#register_form").validate({
        //     rules: {
        //         register_name:{
        //             required: true,
        //         },
        //         register_description: {
        //             required: true,
        //         },
        //         register_price: {
        //             required: true,
        //         }
        //     },
        //     messages: {
        //         register_name: "Please enter a name",
        //         register_description: "Please enter a desctiption",
        //         register_price: "Please enter a price",
        //     },
        //     submitHandler: function(form){
        //         $.ajax({
        //             type:"POST",
        //             url : '/addUser',
        //             data:{
        //                 email : $('#register_email').val(),
        //                 password : $('#register_password').val(),
        //             },
        //             dataType : 'JSON',
        //             success : function(response){
        //                 if(response.status == 'success'){
        //                     $('#registerModal').modal('hide');
        //                 }
        //                console.log("success");
        //             },
        //             error: function(response){
                       
        //                 console.log(response);
        //             }
        //         });
        //     }
        // })    
    })
})