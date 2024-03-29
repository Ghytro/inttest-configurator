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
 * The GithubComGhytroInttestConfiguratorInternalApiAuthCreateUserResponse model module.
 * @module model/GithubComGhytroInttestConfiguratorInternalApiAuthCreateUserResponse
 * @version 2.0
 */
class GithubComGhytroInttestConfiguratorInternalApiAuthCreateUserResponse {
    /**
     * Constructs a new <code>GithubComGhytroInttestConfiguratorInternalApiAuthCreateUserResponse</code>.
     * @alias module:model/GithubComGhytroInttestConfiguratorInternalApiAuthCreateUserResponse
     */
    constructor() { 
        
        GithubComGhytroInttestConfiguratorInternalApiAuthCreateUserResponse.initialize(this);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj) { 
    }

    /**
     * Constructs a <code>GithubComGhytroInttestConfiguratorInternalApiAuthCreateUserResponse</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/GithubComGhytroInttestConfiguratorInternalApiAuthCreateUserResponse} obj Optional instance to populate.
     * @return {module:model/GithubComGhytroInttestConfiguratorInternalApiAuthCreateUserResponse} The populated <code>GithubComGhytroInttestConfiguratorInternalApiAuthCreateUserResponse</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new GithubComGhytroInttestConfiguratorInternalApiAuthCreateUserResponse();

            if (data.hasOwnProperty('created_at')) {
                obj['created_at'] = ApiClient.convertToType(data['created_at'], 'String');
            }
            if (data.hasOwnProperty('id')) {
                obj['id'] = ApiClient.convertToType(data['id'], 'Number');
            }
        }
        return obj;
    }


}

/**
 * @member {String} created_at
 */
GithubComGhytroInttestConfiguratorInternalApiAuthCreateUserResponse.prototype['created_at'] = undefined;

/**
 * @member {Number} id
 */
GithubComGhytroInttestConfiguratorInternalApiAuthCreateUserResponse.prototype['id'] = undefined;






export default GithubComGhytroInttestConfiguratorInternalApiAuthCreateUserResponse;

