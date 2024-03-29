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
 * The ConfiguratorInternalApiAuthCreateUserRequest model module.
 * @module model/ConfiguratorInternalApiAuthCreateUserRequest
 * @version 2.0
 */
class ConfiguratorInternalApiAuthCreateUserRequest {
    /**
     * Constructs a new <code>ConfiguratorInternalApiAuthCreateUserRequest</code>.
     * @alias module:model/ConfiguratorInternalApiAuthCreateUserRequest
     */
    constructor() { 
        
        ConfiguratorInternalApiAuthCreateUserRequest.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>ConfiguratorInternalApiAuthCreateUserRequest</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/ConfiguratorInternalApiAuthCreateUserRequest} obj Optional instance to populate.
     * @return {module:model/ConfiguratorInternalApiAuthCreateUserRequest} The populated <code>ConfiguratorInternalApiAuthCreateUserRequest</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new ConfiguratorInternalApiAuthCreateUserRequest();

            if (data.hasOwnProperty('password')) {
                obj['password'] = ApiClient.convertToType(data['password'], 'String');
            }
            if (data.hasOwnProperty('roles')) {
                obj['roles'] = ApiClient.convertToType(data['roles'], ['String']);
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
ConfiguratorInternalApiAuthCreateUserRequest.prototype['password'] = undefined;

/**
 * @member {Array.<String>} roles
 */
ConfiguratorInternalApiAuthCreateUserRequest.prototype['roles'] = undefined;

/**
 * @member {String} username
 */
ConfiguratorInternalApiAuthCreateUserRequest.prototype['username'] = undefined;






export default ConfiguratorInternalApiAuthCreateUserRequest;

