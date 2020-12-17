"""Script to delete resources in FHIR."""
import json
import logging
import os

from google.auth.transport import requests
from google.oauth2 import service_account

BASE_URL = os.environ["BASE_URL"]
PROJECT_ID = os.environ["PROJECT_ID"]
CLOUD_REGION = os.environ["CLOUD_REGION"]
DATASET_ID = os.environ["DATASET_ID"]
FHIR_STORE_ID = os.environ["FHIR_STORE_ID"]
RESOURCE_TYPE = os.environ["RESOURCE_TYPE"]
URL = f"{BASE_URL}/projects/{PROJECT_ID}/locations/{CLOUD_REGION}"

logging.basicConfig()
logging.getLogger().setLevel(logging.DEBUG)


def get_session():
    """Create an authorized requests session."""
    credentials = service_account.Credentials.from_service_account_file(
        os.environ["GOOGLE_APPLICATION_CREDENTIALS"]
    )
    scoped_credentials = credentials.with_scopes([os.environ["AUTH_URL"]])
    session = requests.AuthorizedSession(scoped_credentials)

    return session


def get_patient_everything(patient_id):
    """Gets all the resources in the patient compartment."""
    resource_path = (
        f"{URL}/datasets/{DATASET_ID}/fhirStores/"
        f"{FHIR_STORE_ID}/fhir/Patient/{patient_id}"
    )
    resource_path += "/$everything"

    session = get_session()

    headers = {"Content-Type": "application/fhir+json;charset=utf-8"}
    response = session.get(resource_path, headers=headers)
    response.raise_for_status()
    resource = response.json()

    return resource, session


def delete_resource(patient_id):
    """
    Delete a FHIR resource.

    Regardless of whether the operation succeeds or
    fails, the server returns a 200 OK HTTP status code. To check that the
    resource was successfully deleted, search for or get the resource and
    see if it exists.
    """
    resource, session = get_patient_everything(patient_id)
    entires = resource["entry"]
    for entry in entires:
        resource_id = entry["resource"]["id"]
        resource_type = entry["resource"]["resourceType"]
        delete_resource_helper(resource_id, resource_type, session)
    return


def delete_resource_helper(resource_id, resource_type, session):
    """Delete."""
    resource_path = (
        f"{URL}/datasets/{DATASET_ID}/fhirStores/"
        f"{FHIR_STORE_ID}/fhir/{resource_type}/{resource_id}"
    )

    response = session.delete(resource_path)
    if response.status_code != 200:
        logging.info(f"Status code raised: {response.status_code}")
        logging.debug(f"{response.json()}")
        return

    logging.info("All patient records have been deleted.")


def get_patient_resources():
    """Get a FHIR resource (Patient for this case)."""
    resource_path = (
        f"{URL}/datasets/{DATASET_ID}/fhirStores/"
        f"{FHIR_STORE_ID}/fhir/Patient"
    )

    session = get_session()

    headers = {"Content-Type": "application/fhir+json;charset=utf-8"}
    response = session.get(resource_path, headers=headers)
    response.raise_for_status()

    if response.status_code != 200:
        return

    resource = response.json()

    return resource


def delete_all_resources():
    """
    Delete a FHIR resource.

    Regardless of whether the operation succeeds or
    fails, the server returns a 200 OK HTTP status code. To check that the
    resource was successfully deleted, search for or get the resource and
    see if it exists.
    """
    patient_ids = []
    resource = get_patient_resources()
    try:
        entires = resource["entry"]
        for entry in entires:
            patient_ids.append(entry["resource"]["id"])

        count = len(patient_ids)
        deleted_resource = 0
        for patient_id in patient_ids:
            while count != delete_resource:
                delete_resource(patient_id)
                deleted_resource += 1

        logging.info(f"{deleted_resource} out of {count} have been deleted.")

    except KeyError:
        return


if __name__ == "__main__":
    print(
        "\033[1;37;40mPlease verify these variables from your environment before proceeding:"
    )
    print(f"\033[1;32;40m \tDATASET_ID: \t {DATASET_ID}")
    print(f"\033[1;32;40m \tFHIR_STORE_ID \t {FHIR_STORE_ID}\n")
    try:
        patient_id = input(
            "\033[1;37;40mEnter a patient id to delete a single patient record "
            "\033[1;31;40m(pressing enter and leaving it blank will DELETE ALL \N{bomb} \N{bomb} the patient records): \n"
        )
        yes_options = ["YES", "Yes", "yes", "Y", "y"]

        if patient_id == "":
            warning_input = input(
                "\033[1;31;40mAre you sure that you want to DELETE ALL \N{bomb} \N{bomb} the patients records? [Y/N]: \n"
            )
            if warning_input in yes_options:
                delete_all_resources()
            else:
                logging.info("\033[1;32;40mThank you for confirming")

        else:
            inp = input(
                f"\033[1;37;40mPlease confirm you want to delete the Patient data with ID {patient_id} "
                f"from {URL}/datasets/{DATASET_ID}/fhirStores/"
                f"{FHIR_STORE_ID}/fhir/ [Y/N] \n"
            )
            if inp in yes_options:
                delete_resource(patient_id)
            else:
                logging.info("\033[1;32;40mThank you for confirming")

    except KeyboardInterrupt:
        logging.info("\033[1;32;40mExiting gracefully.")
