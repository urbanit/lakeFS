/*
 * lakeFS API
 * lakeFS HTTP API
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package io.lakefs.clients.api.model;

import java.util.Objects;
import java.util.Arrays;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

/**
 * PresignMultipartUpload
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen")
public class PresignMultipartUpload {
  public static final String SERIALIZED_NAME_UPLOAD_ID = "upload_id";
  @SerializedName(SERIALIZED_NAME_UPLOAD_ID)
  private String uploadId;

  public static final String SERIALIZED_NAME_PHYSICAL_ADDRESS = "physical_address";
  @SerializedName(SERIALIZED_NAME_PHYSICAL_ADDRESS)
  private String physicalAddress;

  public static final String SERIALIZED_NAME_PRESIGNED_URLS = "presigned_urls";
  @SerializedName(SERIALIZED_NAME_PRESIGNED_URLS)
  private List<String> presignedUrls = null;


  public PresignMultipartUpload uploadId(String uploadId) {
    
    this.uploadId = uploadId;
    return this;
  }

   /**
   * Get uploadId
   * @return uploadId
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "")

  public String getUploadId() {
    return uploadId;
  }


  public void setUploadId(String uploadId) {
    this.uploadId = uploadId;
  }


  public PresignMultipartUpload physicalAddress(String physicalAddress) {
    
    this.physicalAddress = physicalAddress;
    return this;
  }

   /**
   * Get physicalAddress
   * @return physicalAddress
  **/
  @javax.annotation.Nonnull
  @ApiModelProperty(required = true, value = "")

  public String getPhysicalAddress() {
    return physicalAddress;
  }


  public void setPhysicalAddress(String physicalAddress) {
    this.physicalAddress = physicalAddress;
  }


  public PresignMultipartUpload presignedUrls(List<String> presignedUrls) {
    
    this.presignedUrls = presignedUrls;
    return this;
  }

  public PresignMultipartUpload addPresignedUrlsItem(String presignedUrlsItem) {
    if (this.presignedUrls == null) {
      this.presignedUrls = new ArrayList<String>();
    }
    this.presignedUrls.add(presignedUrlsItem);
    return this;
  }

   /**
   * Get presignedUrls
   * @return presignedUrls
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public List<String> getPresignedUrls() {
    return presignedUrls;
  }


  public void setPresignedUrls(List<String> presignedUrls) {
    this.presignedUrls = presignedUrls;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    PresignMultipartUpload presignMultipartUpload = (PresignMultipartUpload) o;
    return Objects.equals(this.uploadId, presignMultipartUpload.uploadId) &&
        Objects.equals(this.physicalAddress, presignMultipartUpload.physicalAddress) &&
        Objects.equals(this.presignedUrls, presignMultipartUpload.presignedUrls);
  }

  @Override
  public int hashCode() {
    return Objects.hash(uploadId, physicalAddress, presignedUrls);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class PresignMultipartUpload {\n");
    sb.append("    uploadId: ").append(toIndentedString(uploadId)).append("\n");
    sb.append("    physicalAddress: ").append(toIndentedString(physicalAddress)).append("\n");
    sb.append("    presignedUrls: ").append(toIndentedString(presignedUrls)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}

