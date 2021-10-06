# coding: utf-8

"""
    Kubeflow Training SDK

    Python SDK for Kubeflow Training  # noqa: E501

    The version of the OpenAPI document: v1.3.0
    Generated by: https://openapi-generator.tech
"""


import pprint
import re  # noqa: F401

import six

from kubeflow.training.configuration import Configuration


class V1JobStatus(object):
    """NOTE: This class is auto generated by OpenAPI Generator.
    Ref: https://openapi-generator.tech

    Do not edit the class manually.
    """

    """
    Attributes:
      openapi_types (dict): The key is attribute name
                            and the value is attribute type.
      attribute_map (dict): The key is attribute name
                            and the value is json key in definition.
    """
    openapi_types = {
        'completion_time': 'V1Time',
        'conditions': 'list[V1JobCondition]',
        'last_reconcile_time': 'V1Time',
        'replica_statuses': 'dict(str, V1ReplicaStatus)',
        'start_time': 'V1Time'
    }

    attribute_map = {
        'completion_time': 'completionTime',
        'conditions': 'conditions',
        'last_reconcile_time': 'lastReconcileTime',
        'replica_statuses': 'replicaStatuses',
        'start_time': 'startTime'
    }

    def __init__(self, completion_time=None, conditions=None, last_reconcile_time=None, replica_statuses=None, start_time=None, local_vars_configuration=None):  # noqa: E501
        """V1JobStatus - a model defined in OpenAPI"""  # noqa: E501
        if local_vars_configuration is None:
            local_vars_configuration = Configuration()
        self.local_vars_configuration = local_vars_configuration

        self._completion_time = None
        self._conditions = None
        self._last_reconcile_time = None
        self._replica_statuses = None
        self._start_time = None
        self.discriminator = None

        if completion_time is not None:
            self.completion_time = completion_time
        self.conditions = conditions
        if last_reconcile_time is not None:
            self.last_reconcile_time = last_reconcile_time
        self.replica_statuses = replica_statuses
        if start_time is not None:
            self.start_time = start_time

    @property
    def completion_time(self):
        """Gets the completion_time of this V1JobStatus.  # noqa: E501


        :return: The completion_time of this V1JobStatus.  # noqa: E501
        :rtype: V1Time
        """
        return self._completion_time

    @completion_time.setter
    def completion_time(self, completion_time):
        """Sets the completion_time of this V1JobStatus.


        :param completion_time: The completion_time of this V1JobStatus.  # noqa: E501
        :type: V1Time
        """

        self._completion_time = completion_time

    @property
    def conditions(self):
        """Gets the conditions of this V1JobStatus.  # noqa: E501

        Conditions is an array of current observed job conditions.  # noqa: E501

        :return: The conditions of this V1JobStatus.  # noqa: E501
        :rtype: list[V1JobCondition]
        """
        return self._conditions

    @conditions.setter
    def conditions(self, conditions):
        """Sets the conditions of this V1JobStatus.

        Conditions is an array of current observed job conditions.  # noqa: E501

        :param conditions: The conditions of this V1JobStatus.  # noqa: E501
        :type: list[V1JobCondition]
        """
        if self.local_vars_configuration.client_side_validation and conditions is None:  # noqa: E501
            raise ValueError("Invalid value for `conditions`, must not be `None`")  # noqa: E501

        self._conditions = conditions

    @property
    def last_reconcile_time(self):
        """Gets the last_reconcile_time of this V1JobStatus.  # noqa: E501


        :return: The last_reconcile_time of this V1JobStatus.  # noqa: E501
        :rtype: V1Time
        """
        return self._last_reconcile_time

    @last_reconcile_time.setter
    def last_reconcile_time(self, last_reconcile_time):
        """Sets the last_reconcile_time of this V1JobStatus.


        :param last_reconcile_time: The last_reconcile_time of this V1JobStatus.  # noqa: E501
        :type: V1Time
        """

        self._last_reconcile_time = last_reconcile_time

    @property
    def replica_statuses(self):
        """Gets the replica_statuses of this V1JobStatus.  # noqa: E501

        ReplicaStatuses is map of ReplicaType and ReplicaStatus, specifies the status of each replica.  # noqa: E501

        :return: The replica_statuses of this V1JobStatus.  # noqa: E501
        :rtype: dict(str, V1ReplicaStatus)
        """
        return self._replica_statuses

    @replica_statuses.setter
    def replica_statuses(self, replica_statuses):
        """Sets the replica_statuses of this V1JobStatus.

        ReplicaStatuses is map of ReplicaType and ReplicaStatus, specifies the status of each replica.  # noqa: E501

        :param replica_statuses: The replica_statuses of this V1JobStatus.  # noqa: E501
        :type: dict(str, V1ReplicaStatus)
        """
        if self.local_vars_configuration.client_side_validation and replica_statuses is None:  # noqa: E501
            raise ValueError("Invalid value for `replica_statuses`, must not be `None`")  # noqa: E501

        self._replica_statuses = replica_statuses

    @property
    def start_time(self):
        """Gets the start_time of this V1JobStatus.  # noqa: E501


        :return: The start_time of this V1JobStatus.  # noqa: E501
        :rtype: V1Time
        """
        return self._start_time

    @start_time.setter
    def start_time(self, start_time):
        """Sets the start_time of this V1JobStatus.


        :param start_time: The start_time of this V1JobStatus.  # noqa: E501
        :type: V1Time
        """

        self._start_time = start_time

    def to_dict(self):
        """Returns the model properties as a dict"""
        result = {}

        for attr, _ in six.iteritems(self.openapi_types):
            value = getattr(self, attr)
            if isinstance(value, list):
                result[attr] = list(map(
                    lambda x: x.to_dict() if hasattr(x, "to_dict") else x,
                    value
                ))
            elif hasattr(value, "to_dict"):
                result[attr] = value.to_dict()
            elif isinstance(value, dict):
                result[attr] = dict(map(
                    lambda item: (item[0], item[1].to_dict())
                    if hasattr(item[1], "to_dict") else item,
                    value.items()
                ))
            else:
                result[attr] = value

        return result

    def to_str(self):
        """Returns the string representation of the model"""
        return pprint.pformat(self.to_dict())

    def __repr__(self):
        """For `print` and `pprint`"""
        return self.to_str()

    def __eq__(self, other):
        """Returns true if both objects are equal"""
        if not isinstance(other, V1JobStatus):
            return False

        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        if not isinstance(other, V1JobStatus):
            return True

        return self.to_dict() != other.to_dict()
