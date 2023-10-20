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

    
    $('.edit').on('click', function(){
      id = $(this).data('id');
      $.ajax
      ({
        type:"GET",
        url : '/products/' + id,
        dataType : 'JSON',
        success : function(response){
            $('#update_name').val(response.Name)
            $('#update_description').val(response.Description)
            $('#update_price').val(response.Price)
            $('#update_id').val(response.ID)
            $('#updateModal').modal('show');
        },
        error: function(response){
            console.log(response);
        }
      });
    })

    $('#update').on('click', function(){
        $('#update_form').attr('action', $('#update_form').attr('action') + "/" + $('#update_id').val());
        $('#update_form').submit()
    })

    $('.delete').on('click', function(){
        id = $(this).data('id');
        bootbox.confirm({
            size: 'small',
            message: 'Are you sure?',
            callback: function(result) {
                 /* result is a boolean; true = OK, false = Cancel*/ 
                    if (result){
                        $('#delete_form').attr('action', $('#delete_form').attr('action') + "/" + id);
                        $('#delete_form').submit()
                    }
                }
            });
    })
})