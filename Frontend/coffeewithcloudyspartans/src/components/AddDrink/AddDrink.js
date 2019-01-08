import React from 'react';
import { Container, Row, Col, Button, Fa, Card, CardBody, ModalFooter } from 'mdbreact';
import {Redirect} from 'react-router-dom';
import cookie from 'react-cookies';
import DrinkPhotos from '../DrinkPhotos/DrinkPhotos';
import axios from 'axios';
var swal = require('sweetalert')

const URL = "http://localhost:4004"

class AddDrink extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      DrinkName : '',
      Price : '',
      Size : '',
      Description: '',
      Redirection_Value : false,
      errors : false,
      isUploaded : false,
      showPhotos : false,
      photos : []
    }
  }

  onDrop = (acceptedFiles, rejectedFiles) => {
    var photos = [...this.state.photos]
     photos.push(acceptedFiles);
     this.setState({
       photos: photos,
       isUploaded: true,
       showPhotos: true
     });
     console.log("photos", photos);
 }
 
 handleDeleteFile = (event, name) => {
     event.preventDefault();
     console.log("event.target.value", name);
     var fileName = name;
     var oldPhotos = [...this.state.photos];
 
     var newPhotos = [];
     for (let position = 0; position < oldPhotos.length; position++) {
       if (oldPhotos[position][0].name === fileName) {
         newPhotos = oldPhotos.splice(position, 1);
       }
     }
     console.log("newPhotos", newPhotos);
 
     {
       this.setState({
         photos: oldPhotos
       });
     }
   };


uploadFiles = files => {
    const uploadFiles = new FormData();
    var filenames = "";
    for (let index = 0; index < files.length; index++) {
      if (index === files.length - 1) {
        filenames = filenames.concat(files[index][0].name);
      } else {
        filenames = filenames.concat(files[index][0].name + ",");
      }
      uploadFiles.append("file", files[index][0]);
    }

    axios.post(URL+ "/upload-files",uploadFiles).then(
      response => {
        console.log(response.data.message);
      },
      error => {
        console.log(error);
      }
    );

    return filenames;
  };

  onDrinkNameChange = (event) => {
    this.setState({DrinkName: event.target.value})
  }

  onPriceChange = (event) => {
    this.setState({Price: event.target.value})
  }

  onSizeChange = (event) => {
    this.setState({Size: event.target.value})
  }

   onDescriptionChange = (event) => {
    this.setState({Description: event.target.value})
  }


  onSubmitSignIn = () => {
    let file = this.state.photos;
    if (file.length > 0) {
      let filenames = "";
      
      fetch(URL+'/addfolder', {
        method: 'post',
        headers: {'Content-Type': 'application/json'},
        credentials : 'include',
        body: JSON.stringify({
            FolderName: this.state.DrinkName,
        })  
      }).
      then(response => {
      filenames = this.uploadFiles(file)
      fetch(URL + '/addadrink', {
        method: 'post',
        headers: {'Content-Type': 'application/json'},
        credentials : 'include',
        body: JSON.stringify({
            Name: this.state.DrinkName,
            Price: this.state.Price,
            Size: this.state.Size,
            Description: this.state.Description
    })
  })
  .then(response => {
    if(response.status === 400)
      {
        this.setState({errors : true})
      }
    else
      {
        if(response.status === 200)
        {
          swal("Drink Added successfully!", " ", "success");
            this.setState({Redirecttohome : true})
        }
     }
        })
      })    
}
else
{
    window.alert("Uploading Photos is Mandatory");
}
}


  render()
  {
    let Redirecty = null;
    let Errors = null;
   
      if(this.state.errors === true)
    {
      Errors = <p class="error">Error Adding Drink</p>
    }
    var Redirecttohome = null;
    if(this.state.Redirecttohome === true){
        Redirecttohome = (<Redirect to="/home" />)
    }

    var USERTYPE = localStorage.getItem("ACCOUNTTYPE")
    var Redirecttologin = null;
    if(USERTYPE!="admin"){
      Redirecttologin = (<Redirect to="/admin/login"/>)
    }


    return(
      <div>
      {Redirecty}
      {Redirecttologin}
      {Redirecttohome}
      <Container>
       <br />
       <h1 class="makeitcenetersignup">Add a Drink to Coffee with Cloudy Spartans </h1>
       <br />
       <br />
        <section className="form-elegant">
          <Row >
            <Col md="4" className="mx-auto">
              <Card>
                <CardBody className="mx-4">
                  <div className="text-center">
                    <hr></hr>
                  </div>
                  <input type="text" class="form-control" id="exampleInputFirstName" aria-describedby="emailHelp" placeholder="Name of the drink"  onChange={this.onDrinkNameChange} required/>
                  <br>
                  </br>
                  <input type="text" class="form-control" id="exampleInputSecondName" aria-describedby="emailHelp" placeholder="Price of the drink"  onChange={this.onPriceChange} required/>
                  <br>
                  </br>
                  <input type="email" class="form-control" id="exampleInputEmail1" aria-describedby="emailHelp" placeholder="Size of the drink"  onChange={this.onSizeChange} required/>
                  <br>
                  </br>
                  <input type="text" class="form-control" id="exampleInputusername" placeholder="Description" onChange={this.onDescriptionChange} required/>
                   {<DrinkPhotos 
                        photos = {this.state.photos}
                        isUploaded = {this.state.isUploaded}
                        showPhotos = {this.state.showPhotos}
                        onDrop = {this.onDrop}
                        handleDeleteFile = {this.handleDeleteFile}
                     />}    
                  <div className="text-center mb-3">
                    <Button type="button" gradient="blue" className="btn btn-primary btn-lg btn-block" onClick = {this.onSubmitSignIn}>Add</Button>
                      <hr></hr>
                  </div>
                </CardBody>     
              </Card>
            </Col>
          </Row>
        </section>
         {Errors}
      </Container>
      </div>
      );
  }
}

export default AddDrink;