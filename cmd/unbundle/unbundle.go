package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/fhir/go/fhirversion"
	"github.com/google/fhir/go/jsonformat"
	datatypes_go_proto "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
	rpb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	patient_go_proto "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/patient_go_proto"
	"google.golang.org/protobuf/proto"
)

// Unbundle takes a filename and unbundles the resources into separate files.
func Unbundle(filename, outputDir string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed opening file: %s", err)
	}

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// get current tz as a string:
	tz, _ := time.Now().Zone()
	tz = "UTC"
	um, err := jsonformat.NewUnmarshaller(tz, fhirversion.R4)
	if err != nil {
		return fmt.Errorf("failed to create unmarshaller: %v", err)
	}

	o, err := um.Unmarshal(byteValue)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %v", err)
	}

	switch r := o.(type) {
	case *rpb.ContainedResource:
		log.Printf("Unbundled resource: %T", r)
		log.Printf("Unbundled oneof: %T", r.OneofResource)
		switch rr := r.OneofResource.(type) {
		case *rpb.ContainedResource_Bundle:
			log.Printf("bundle: %T", rr)
			for i, e := range rr.Bundle.Entry {
				if err := writeBundleEntry(filename, outputDir, i, e); err != nil {
					return fmt.Errorf("failed to write bundle entry: %v", err)
				}
			}
		default:
			log.Printf("not a bundle: %T", o)
		}
	}

	log.Printf("Successfully unbundled resource: %T", o)
	return nil
}

type Resource interface {
	GetId() *datatypes_go_proto.Id
	GetMeta() *datatypes_go_proto.Meta
}

var p = &patient_go_proto.Patient{}

func writeBundleEntry(filename, outputDir string, index int, e *rpb.Bundle_Entry) error {
	// {fname}-{resourceType}-{index}-{id}.json
	rr, err := getResource(e.Resource)
	if err != nil {
		return fmt.Errorf("failed to get resource: %v", err)
	}
	r, ok := rr.(Resource)
	if !ok {
		fmt.Printf("not a resource: %T\n", rr)
		return nil
	}
	ts := fmt.Sprintf("%T", rr)
	_, t, _ := strings.Cut(ts, ".")
	// filename with extension:
	fn := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	fname := fmt.Sprintf("%s-%s-%d-%s.json", fn, t, index, r.GetId().GetValue())
	path := filepath.Join(outputDir, fname)

	// create file:
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer f.Close()

	// write to file:
	m, err := jsonformat.NewPrettyMarshaller(fhirversion.R4)
	if err != nil {
		return fmt.Errorf("failed to create marshaller: %v", err)
	}
	b, err := m.Marshal(e.Resource)
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}
	if _, err := f.Write(b); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}
	return nil
}

func getResource(r *rpb.ContainedResource) (proto.Message, error) {
	switch r := r.OneofResource.(type) {
	case *rpb.ContainedResource_Account:
		return r.Account, nil
	case *rpb.ContainedResource_ActivityDefinition:
		return r.ActivityDefinition, nil
	case *rpb.ContainedResource_AdverseEvent:
		return r.AdverseEvent, nil
	case *rpb.ContainedResource_AllergyIntolerance:
		return r.AllergyIntolerance, nil
	case *rpb.ContainedResource_Appointment:
		return r.Appointment, nil
	case *rpb.ContainedResource_AppointmentResponse:
		return r.AppointmentResponse, nil
	case *rpb.ContainedResource_AuditEvent:
		return r.AuditEvent, nil
	case *rpb.ContainedResource_Basic:
		return r.Basic, nil
	case *rpb.ContainedResource_Binary:
		return r.Binary, nil
	case *rpb.ContainedResource_BiologicallyDerivedProduct:
		return r.BiologicallyDerivedProduct, nil
	case *rpb.ContainedResource_BodyStructure:
		return r.BodyStructure, nil
	case *rpb.ContainedResource_Bundle:
		return r.Bundle, nil
	case *rpb.ContainedResource_CapabilityStatement:
		return r.CapabilityStatement, nil
	case *rpb.ContainedResource_CarePlan:
		return r.CarePlan, nil
	case *rpb.ContainedResource_CareTeam:
		return r.CareTeam, nil
	case *rpb.ContainedResource_CatalogEntry:
		return r.CatalogEntry, nil
	case *rpb.ContainedResource_ChargeItem:
		return r.ChargeItem, nil
	case *rpb.ContainedResource_ChargeItemDefinition:
		return r.ChargeItemDefinition, nil
	case *rpb.ContainedResource_Claim:
		return r.Claim, nil
	case *rpb.ContainedResource_ClaimResponse:
		return r.ClaimResponse, nil
	case *rpb.ContainedResource_ClinicalImpression:
		return r.ClinicalImpression, nil
	case *rpb.ContainedResource_CodeSystem:
		return r.CodeSystem, nil
	case *rpb.ContainedResource_Communication:
		return r.Communication, nil
	case *rpb.ContainedResource_CommunicationRequest:
		return r.CommunicationRequest, nil
	case *rpb.ContainedResource_CompartmentDefinition:
		return r.CompartmentDefinition, nil
	case *rpb.ContainedResource_Composition:
		return r.Composition, nil
	case *rpb.ContainedResource_ConceptMap:
		return r.ConceptMap, nil
	case *rpb.ContainedResource_Condition:
		return r.Condition, nil
	case *rpb.ContainedResource_Consent:
		return r.Consent, nil
	case *rpb.ContainedResource_Contract:
		return r.Contract, nil
	case *rpb.ContainedResource_Coverage:
		return r.Coverage, nil
	case *rpb.ContainedResource_CoverageEligibilityRequest:
		return r.CoverageEligibilityRequest, nil
	case *rpb.ContainedResource_CoverageEligibilityResponse:
		return r.CoverageEligibilityResponse, nil
	case *rpb.ContainedResource_DetectedIssue:
		return r.DetectedIssue, nil
	case *rpb.ContainedResource_Device:
		return r.Device, nil
	case *rpb.ContainedResource_DeviceDefinition:
		return r.DeviceDefinition, nil
	case *rpb.ContainedResource_DeviceMetric:
		return r.DeviceMetric, nil
	case *rpb.ContainedResource_DeviceRequest:
		return r.DeviceRequest, nil
	case *rpb.ContainedResource_DeviceUseStatement:
		return r.DeviceUseStatement, nil
	case *rpb.ContainedResource_DiagnosticReport:
		return r.DiagnosticReport, nil
	case *rpb.ContainedResource_DocumentManifest:
		return r.DocumentManifest, nil
	case *rpb.ContainedResource_DocumentReference:
		return r.DocumentReference, nil
	case *rpb.ContainedResource_EffectEvidenceSynthesis:
		return r.EffectEvidenceSynthesis, nil
	case *rpb.ContainedResource_Encounter:
		return r.Encounter, nil
	case *rpb.ContainedResource_Endpoint:
		return r.Endpoint, nil
	case *rpb.ContainedResource_EnrollmentRequest:
		return r.EnrollmentRequest, nil
	case *rpb.ContainedResource_EnrollmentResponse:
		return r.EnrollmentResponse, nil
	case *rpb.ContainedResource_EpisodeOfCare:
		return r.EpisodeOfCare, nil
	case *rpb.ContainedResource_EventDefinition:
		return r.EventDefinition, nil
	case *rpb.ContainedResource_Evidence:
		return r.Evidence, nil
	case *rpb.ContainedResource_EvidenceVariable:
		return r.EvidenceVariable, nil
	case *rpb.ContainedResource_ExampleScenario:
		return r.ExampleScenario, nil
	case *rpb.ContainedResource_ExplanationOfBenefit:
		return r.ExplanationOfBenefit, nil
	case *rpb.ContainedResource_FamilyMemberHistory:
		return r.FamilyMemberHistory, nil
	case *rpb.ContainedResource_Flag:
		return r.Flag, nil
	case *rpb.ContainedResource_Goal:
		return r.Goal, nil
	case *rpb.ContainedResource_GraphDefinition:
		return r.GraphDefinition, nil
	case *rpb.ContainedResource_Group:
		return r.Group, nil
	case *rpb.ContainedResource_GuidanceResponse:
		return r.GuidanceResponse, nil
	case *rpb.ContainedResource_HealthcareService:
		return r.HealthcareService, nil
	case *rpb.ContainedResource_ImagingStudy:
		return r.ImagingStudy, nil
	case *rpb.ContainedResource_Immunization:
		return r.Immunization, nil
	case *rpb.ContainedResource_ImmunizationEvaluation:
		return r.ImmunizationEvaluation, nil
	case *rpb.ContainedResource_ImmunizationRecommendation:
		return r.ImmunizationRecommendation, nil
	case *rpb.ContainedResource_ImplementationGuide:
		return r.ImplementationGuide, nil
	case *rpb.ContainedResource_InsurancePlan:
		return r.InsurancePlan, nil
	case *rpb.ContainedResource_Invoice:
		return r.Invoice, nil
	case *rpb.ContainedResource_Library:
		return r.Library, nil
	case *rpb.ContainedResource_Linkage:
		return r.Linkage, nil
	case *rpb.ContainedResource_List:
		return r.List, nil
	case *rpb.ContainedResource_Location:
		return r.Location, nil
	case *rpb.ContainedResource_Measure:
		return r.Measure, nil
	case *rpb.ContainedResource_MeasureReport:
		return r.MeasureReport, nil
	case *rpb.ContainedResource_Media:
		return r.Media, nil
	case *rpb.ContainedResource_Medication:
		return r.Medication, nil
	case *rpb.ContainedResource_MedicationAdministration:
		return r.MedicationAdministration, nil
	case *rpb.ContainedResource_MedicationDispense:
		return r.MedicationDispense, nil
	case *rpb.ContainedResource_MedicationKnowledge:
		return r.MedicationKnowledge, nil
	case *rpb.ContainedResource_MedicationRequest:
		return r.MedicationRequest, nil
	case *rpb.ContainedResource_MedicationStatement:
		return r.MedicationStatement, nil
	case *rpb.ContainedResource_MedicinalProduct:
		return r.MedicinalProduct, nil
	case *rpb.ContainedResource_MedicinalProductAuthorization:
		return r.MedicinalProductAuthorization, nil
	case *rpb.ContainedResource_MedicinalProductContraindication:
		return r.MedicinalProductContraindication, nil
	case *rpb.ContainedResource_MedicinalProductIndication:
		return r.MedicinalProductIndication, nil
	case *rpb.ContainedResource_MedicinalProductIngredient:
		return r.MedicinalProductIngredient, nil
	case *rpb.ContainedResource_MedicinalProductInteraction:
		return r.MedicinalProductInteraction, nil
	case *rpb.ContainedResource_MedicinalProductManufactured:
		return r.MedicinalProductManufactured, nil
	case *rpb.ContainedResource_MedicinalProductPackaged:
		return r.MedicinalProductPackaged, nil
	case *rpb.ContainedResource_MedicinalProductPharmaceutical:
		return r.MedicinalProductPharmaceutical, nil
	case *rpb.ContainedResource_MedicinalProductUndesirableEffect:
		return r.MedicinalProductUndesirableEffect, nil
	case *rpb.ContainedResource_MessageDefinition:
		return r.MessageDefinition, nil
	case *rpb.ContainedResource_MessageHeader:
		return r.MessageHeader, nil
	case *rpb.ContainedResource_MolecularSequence:
		return r.MolecularSequence, nil
	case *rpb.ContainedResource_NamingSystem:
		return r.NamingSystem, nil
	case *rpb.ContainedResource_NutritionOrder:
		return r.NutritionOrder, nil
	case *rpb.ContainedResource_Observation:
		return r.Observation, nil
	case *rpb.ContainedResource_ObservationDefinition:
		return r.ObservationDefinition, nil
	case *rpb.ContainedResource_OperationDefinition:
		return r.OperationDefinition, nil
	case *rpb.ContainedResource_OperationOutcome:
		return r.OperationOutcome, nil
	case *rpb.ContainedResource_Organization:
		return r.Organization, nil
	case *rpb.ContainedResource_OrganizationAffiliation:
		return r.OrganizationAffiliation, nil
	case *rpb.ContainedResource_Parameters:
		return r.Parameters, nil
	case *rpb.ContainedResource_Patient:
		return r.Patient, nil
	case *rpb.ContainedResource_PaymentNotice:
		return r.PaymentNotice, nil
	case *rpb.ContainedResource_PaymentReconciliation:
		return r.PaymentReconciliation, nil
	case *rpb.ContainedResource_Person:
		return r.Person, nil
	case *rpb.ContainedResource_PlanDefinition:
		return r.PlanDefinition, nil
	case *rpb.ContainedResource_Practitioner:
		return r.Practitioner, nil
	case *rpb.ContainedResource_PractitionerRole:
		return r.PractitionerRole, nil
	case *rpb.ContainedResource_Procedure:
		return r.Procedure, nil
	case *rpb.ContainedResource_Provenance:
		return r.Provenance, nil
	case *rpb.ContainedResource_Questionnaire:
		return r.Questionnaire, nil
	case *rpb.ContainedResource_QuestionnaireResponse:
		return r.QuestionnaireResponse, nil
	case *rpb.ContainedResource_RelatedPerson:
		return r.RelatedPerson, nil
	case *rpb.ContainedResource_RequestGroup:
		return r.RequestGroup, nil
	case *rpb.ContainedResource_ResearchDefinition:
		return r.ResearchDefinition, nil
	case *rpb.ContainedResource_ResearchElementDefinition:
		return r.ResearchElementDefinition, nil
	case *rpb.ContainedResource_ResearchStudy:
		return r.ResearchStudy, nil
	case *rpb.ContainedResource_ResearchSubject:
		return r.ResearchSubject, nil
	case *rpb.ContainedResource_RiskAssessment:
		return r.RiskAssessment, nil
	case *rpb.ContainedResource_RiskEvidenceSynthesis:
		return r.RiskEvidenceSynthesis, nil
	case *rpb.ContainedResource_Schedule:
		return r.Schedule, nil
	case *rpb.ContainedResource_SearchParameter:
		return r.SearchParameter, nil
	case *rpb.ContainedResource_ServiceRequest:
		return r.ServiceRequest, nil
	case *rpb.ContainedResource_Slot:
		return r.Slot, nil
	case *rpb.ContainedResource_Specimen:
		return r.Specimen, nil
	case *rpb.ContainedResource_SpecimenDefinition:
		return r.SpecimenDefinition, nil
	case *rpb.ContainedResource_StructureDefinition:
		return r.StructureDefinition, nil
	case *rpb.ContainedResource_StructureMap:
		return r.StructureMap, nil
	case *rpb.ContainedResource_Subscription:
		return r.Subscription, nil
	case *rpb.ContainedResource_Substance:
		return r.Substance, nil
	case *rpb.ContainedResource_SubstanceNucleicAcid:
		return r.SubstanceNucleicAcid, nil
	case *rpb.ContainedResource_SubstancePolymer:
		return r.SubstancePolymer, nil
	case *rpb.ContainedResource_SubstanceProtein:
		return r.SubstanceProtein, nil
	case *rpb.ContainedResource_SubstanceReferenceInformation:
		return r.SubstanceReferenceInformation, nil
	case *rpb.ContainedResource_SubstanceSourceMaterial:
		return r.SubstanceSourceMaterial, nil
	case *rpb.ContainedResource_SubstanceSpecification:
		return r.SubstanceSpecification, nil
	case *rpb.ContainedResource_SupplyDelivery:
		return r.SupplyDelivery, nil
	case *rpb.ContainedResource_SupplyRequest:
		return r.SupplyRequest, nil
	case *rpb.ContainedResource_Task:
		return r.Task, nil
	case *rpb.ContainedResource_TerminologyCapabilities:
		return r.TerminologyCapabilities, nil
	case *rpb.ContainedResource_TestReport:
		return r.TestReport, nil
	case *rpb.ContainedResource_TestScript:
		return r.TestScript, nil
	case *rpb.ContainedResource_ValueSet:
		return r.ValueSet, nil
	case *rpb.ContainedResource_VerificationResult:
		return r.VerificationResult, nil
	case *rpb.ContainedResource_VisionPrescription:
		return r.VisionPrescription, nil
	default:
		return nil, fmt.Errorf("unknown resource type %T", r)
	}
}
