import React,{Component} from 'react';
import './Cart.css';
import {Route} from 'react-router-dom'
import Link from 'react-router-dom/Link';


const URL="http://localhost:4004"
class Cart extends Component {
    constructor(props) {
    super(props);
    this.state = {
      photos : []
    }
  }

  componentDidMount()
  {
        fetch(URL+ '/getDrinkImg', {
          method: 'post',
          headers: {'Content-Type': 'application/json'},
          credentials : 'include',
          body : JSON.stringify({
            id : this.props.name
          })
        })
        .then(response => response.json())
        .then(data => {
          let imageArr = []
        for (let i = 0; i < data.results.length; i++) {
          let imagePreview = 'data:image/jpg;charset=utf-8;base64, ' + data.results[i];
                                imageArr.push(imagePreview);
                                const photoArr = this.state.photos.slice();
                                photoArr[i] = imagePreview;
                                this.setState({
                                    photos: photoArr
                                });
                                console.log('Photo State: ', this.state.photos);
                  }
        })
  }
  render()
  {
    let carousalBlock = this.state.photos.map(function (item, index) {

            return (
                <div className={index == 0 ? "carousel-item active" : "carousel-item"} key={index}>
                    <img className="carousel-img property-display-img" src={item} width="350" height="200" alt="property-image" />
                </div>
            )
        });

        let carousalIndicator = this.state.photos.map(function (item, index) {

            return (                
                    <li data-target="#myCarousel" data-slide-to={index} className={index == 0 ? "active" : ""} key={index}></li>     
            )
        });

      return (
        <div className="row justify-content-center ownercard3 ownercard-2" >
         <div className="col-md-12" >
                  <div className="row" style={{padding:'5px'}} >
                  <div className="col-md-6">
              <div class="image-card">
                <div id="myCarousel" className="carousel slide contain" data-ride="carousel">
                  <ul className="carousel-indicators">
                    {carousalIndicator}
                  </ul>
                  <div className="carousel-inner">
                    {carousalBlock}
                  </div>
                </div>
              </div>
                  </div>

           <div className="col-md-6 content descrip">
            <div class="meta">  
              <span class="cinema"> Name : {this.props.name}</span>
            </div>
            <div class="description">
              <p> Desciption : {"Must try special from Coffee with Cloudy Spartans"}</p>
            </div>
            <div class="extra">
              <div class="ui label">Sizes : {this.props.sizes}</div>
              <div class="ui label">Price : {this.props.price}</div>
              <div class="ui label">Count : {this.props.count}</div>
              <br />
              <button className="btn btn-primary" onClick={drinkinfo=>this.props.drinkadd(this.props.drinkinfo)}>+</button>
          &nbsp;&nbsp;<button className="btn btn-primary" onClick={drinkinfo=>this.props.drinkdelete(this.props.drinkinfo)}>-</button>
          <br />
            </div>
          </div>
        </div>
        </div>
        </div>
      );
    }
}

export default Cart;

