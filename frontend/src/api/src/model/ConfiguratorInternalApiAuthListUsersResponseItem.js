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
 * The ConfiguratorInternalApiAuthListUsersResponseItem model module.
 * @module model/ConfiguratorInternalApiAuthListUsersResponseItem
 * @version 2.0
 */
class ConfiguratorInternalApiAuthListUsersResponseItem {
    /**
     * Constructs a new <code>ConfiguratorInternalApiAuthListUsersResponseItem</code>.
     * @alias module:model/ConfiguratorInternalApiAuthListUsersResponseItem
     */
    constructor() { 
        
        ConfiguratorInternalApiAuthListUsersResponseItem.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>ConfiguratorInternalApiAuthListUsersResponseItem</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/ConfiguratorInternalApiAuthListUsersResponseItem} obj Optional instance to populate.
     * @return {module:model/ConfiguratorInternalApiAuthListUsersResponseItem} The populated <code>ConfiguratorInternalApiAuthListUsersResponseItem</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new ConfiguratorInternalApiAuthListUsersResponseItem();

            if (data.hasOwnProperty('created_at')) {
                obj['created_at'] = ApiClient.convertToType(data['created_at'], 'String');
            }
            if (data.hasOwnProperty('id')) {
                obj['id'] = ApiClient.convertToType(data['id'], 'Number');
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
 * @member {String} created_at
 */
ConfiguratorInternalApiAuthListUsersResponseItem.prototype['created_at'] = undefined;

/**
 * @member {Number} id
 */
ConfiguratorInternalApiAuthListUsersResponseItem.prototype['id'] = undefined;

/**
 * @member {Array.<String>} roles
 */
ConfiguratorInternalApiAuthListUsersResponseItem.prototype['roles'] = undefined;

/**
 * @member {String} username
 */
ConfiguratorInternalApiAuthListUsersResponseItem.prototype['username'] = undefined;






export default ConfiguratorInternalApiAuthListUsersResponseItem;

