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
 * The ConfiguratorInternalApiAuthUpdateRoleRequest model module.
 * @module model/ConfiguratorInternalApiAuthUpdateRoleRequest
 * @version 2.0
 */
class ConfiguratorInternalApiAuthUpdateRoleRequest {
    /**
     * Constructs a new <code>ConfiguratorInternalApiAuthUpdateRoleRequest</code>.
     * @alias module:model/ConfiguratorInternalApiAuthUpdateRoleRequest
     */
    constructor() { 
        
        ConfiguratorInternalApiAuthUpdateRoleRequest.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>ConfiguratorInternalApiAuthUpdateRoleRequest</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/ConfiguratorInternalApiAuthUpdateRoleRequest} obj Optional instance to populate.
     * @return {module:model/ConfiguratorInternalApiAuthUpdateRoleRequest} The populated <code>ConfiguratorInternalApiAuthUpdateRoleRequest</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new ConfiguratorInternalApiAuthUpdateRoleRequest();

            if (data.hasOwnProperty('desc')) {
                obj['desc'] = ApiClient.convertToType(data['desc'], 'String');
            }
            if (data.hasOwnProperty('name')) {
                obj['name'] = ApiClient.convertToType(data['name'], 'String');
            }
            if (data.hasOwnProperty('perm_ids')) {
                obj['perm_ids'] = ApiClient.convertToType(data['perm_ids'], ['Number']);
            }
        }
        return obj;
    }


}

/**
 * @member {String} desc
 */
ConfiguratorInternalApiAuthUpdateRoleRequest.prototype['desc'] = undefined;

/**
 * @member {String} name
 */
ConfiguratorInternalApiAuthUpdateRoleRequest.prototype['name'] = undefined;

/**
 * @member {Array.<Number>} perm_ids
 */
ConfiguratorInternalApiAuthUpdateRoleRequest.prototype['perm_ids'] = undefined;






export default ConfiguratorInternalApiAuthUpdateRoleRequest;

