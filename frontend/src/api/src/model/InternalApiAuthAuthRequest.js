/**
 * IntTest configurator
 * idk what to write here it's just a swagger
 *
 * The version of the OpenAPI document: 2.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 *
 */

import ApiClient from '../ApiClient';

/**
 * The InternalApiAuthAuthRequest model module.
 * @module model/InternalApiAuthAuthRequest
 * @version 2.0
 */
class InternalApiAuthAuthRequest {
    /**
     * Constructs a new <code>InternalApiAuthAuthRequest</code>.
     * @alias module:model/InternalApiAuthAuthRequest
     */
    constructor() { 
        
        InternalApiAuthAuthRequest.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>InternalApiAuthAuthRequest</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/InternalApiAuthAuthRequest} obj Optional instance to populate.
     * @return {module:model/InternalApiAuthAuthRequest} The populated <code>InternalApiAuthAuthRequest</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new InternalApiAuthAuthRequest();

            if (data.hasOwnProperty('password')) {
                obj['password'] = ApiClient.convertToType(data['password'], 'String');
            }
            if (data.hasOwnProperty('username')) {
                obj['username'] = ApiClient.convertToType(data['username'], 'String');
            }
        }
        return obj;
    }


}

/**
 * @member {String} password
 */
InternalApiAuthAuthRequest.prototype['password'] = undefined;

/**
 * @member {String} username
 */
InternalApiAuthAuthRequest.prototype['username'] = undefined;






export default InternalApiAuthAuthRequest;
