import React, { Component } from "react";
import Dropzone from "react-dropzone";
import filesize from "filesize";
import './DrinkPhotos.css'

class DrinkPhotos extends Component{
	constructor(props) {
    	super(props);
  	}

render() {
  return (
  	<div class="shadowingcontainer">
      <div class="shadowingcontainerphotos">
      {console.log("LOADED PHOTOS")}
        <Dropzone
          className="file-upload-area"
          onDrop={files => this.props.onDrop(files)}
        >
          <div className="space" />
          <div className="space" />
          <br />
          <h5 className="photoText">Add photos to the drink</h5>
          <hr className="hr" />

          <div>
            {" "}
            <div className="btn btn-plain btn-file-uploader">
              <button className="btn-upload">Select Photos to upload</button>
            </div>
          </div>
          <div>
            {" "}
          </div>
        </Dropzone>
        {this.props.isUploaded && (
          <table className="table-upload">
            <tbody className="table-upload-body">
              {this.props.photos.map(data => (
                <tr
                  key={this.props.photos.indexOf(data)}
                  className="table-upload-row"
                >
                  <td className="table-upload-row-preview">
                    <img
                      className="preview-image"
                      src={
                        data[0].type === "application/pdf"
                          ? "abcd"
                          : URL.createObjectURL(data[0])
                      }
                      alt=""
                    />
                  </td>
                  <td className="table-upload-row-name">
                    <span>{data[0].name}</span>
                  </td>
                  <td className="table-upload-row-size">
                    {filesize(data[0].size)}
                  </td>
                  <td className="table-upload-row-delete">
                    <button
                      value={data[0].name}
                      className="btn btn-danger btn-small"
                      onClick={event => {
                        this.props.handleDeleteFile(event, data[0].name);
                      }}
                    >
                      <span>Delete</span>
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
      </div>
  );
};
}

export default DrinkPhotos;