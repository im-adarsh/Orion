package newrelic

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Config contains Application and Transaction behavior settings.
// Use NewConfig to create a Config with proper defaults.
type Config struct {
	// AppName is used by New Relic to link data across servers.
	//
	// https://docs.newrelic.com/docs/apm/new-relic-apm/installation-configuration/naming-your-application
	AppName string

	// License is your New Relic license key.
	//
	// https://docs.newrelic.com/docs/accounts-partnerships/accounts/account-setup/license-key
	License string

	// Logger controls go-agent logging.  See log.go.
	Logger Logger

	// Enabled determines whether the agent will communicate with the New
	// Relic servers and spawn goroutines.  Setting this to be false can be
	// useful in testing and staging situations.
	Enabled bool

	// Labels are key value pairs used to roll up applications into specific
	// categories.
	//
	// https://docs.newrelic.com/docs/apm/new-relic-apm/maintenance/labels-categories-organizing-your-apps-servers
	Labels map[string]string

	// HighSecurity guarantees that certain agent settings can not be made
	// more permissive.  This setting must match the corresponding account
	// setting in the New Relic UI.
	//
	// https://docs.newrelic.com/docs/accounts-partnerships/accounts/security/high-security
	HighSecurity bool

	// SecurityPoliciesToken enables security policies if set to a non-empty
	// string.  Only set this if security policies have been enabled on your
	// account.  This cannot be used in conjunction with HighSecurity.
	SecurityPoliciesToken string

	// CustomInsightsEvents controls the behavior of
	// Application.RecordCustomEvent.
	//
	// https://docs.newrelic.com/docs/insights/new-relic-insights/adding-querying-data/inserting-custom-events-new-relic-apm-agents
	CustomInsightsEvents struct {
		// Enabled controls whether RecordCustomEvent will collect
		// custom analytics events.  High security mode overrides this
		// setting.
		Enabled bool
	}

	// TransactionEvents controls the behavior of transaction analytics
	// events.
	TransactionEvents struct {
		// Enabled controls whether transaction events are captured.
		Enabled bool
		// Attributes controls the attributes included with transaction
		// events.
		Attributes AttributeDestinationConfig
	}

	// ErrorCollector controls the capture of errors.
	ErrorCollector struct {
		// Enabled controls whether errors are captured.  This setting
		// affects both traced errors and error analytics events.
		Enabled bool
		// CaptureEvents controls whether error analytics events are
		// captured.
		CaptureEvents bool
		// IgnoreStatusCodes controls which http response codes are
		// automatically turned into errors.  By default, response codes
		// greater than or equal to 400, with the exception of 404, are
		// turned into errors.
		IgnoreStatusCodes []int
		// Attributes controls the attributes included with errors.
		Attributes AttributeDestinationConfig
	}

	// TransactionTracer controls the capture of transaction traces.
	TransactionTracer struct {
		// Enabled controls whether transaction traces are captured.
		Enabled bool
		// Threshold controls whether a transaction trace will be
		// considered for capture.  Of the traces exceeding the
		// threshold, the slowest trace every minute is captured.
		Threshold struct {
			// If IsApdexFailing is true then the trace threshold is
			// four times the apdex threshold.
			IsApdexFailing bool
			// If IsApdexFailing is false then this field is the
			// threshold, otherwise it is ignored.
			Duration time.Duration
		}
		// SegmentThreshold is the threshold at which segments will be
		// added to the trace.  Lowering this setting may increase
		// overhead.
		SegmentThreshold time.Duration
		// StackTraceThreshold is the threshold at which segments will
		// be given a stack trace in the transaction trace.  Lowering
		// this setting will drastically increase overhead.
		StackTraceThreshold time.Duration
		// Attributes controls the attributes included with transaction
		// traces.
		Attributes AttributeDestinationConfig
	}

	// HostDisplayName gives this server a recognizable name in the New
	// Relic UI.  This is an optional setting.
	HostDisplayName string

	// Transport customizes http.Client communication with New Relic
	// servers.  This may be used to configure a proxy.
	Transport http.RoundTripper

	// Utilization controls the detection and gathering of system
	// information.
	Utilization struct {
		// DetectAWS controls whether the Application attempts to detect
		// AWS.
		DetectAWS bool
		// DetectAzure controls whether the Application attempts to detect
		// Azure.
		DetectAzure bool
		// DetectPCF controls whether the Application attempts to detect
		// PCF.
		DetectPCF bool
		// DetectGCP controls whether the Application attempts to detect
		// GCP.
		DetectGCP bool
		// DetectDocker controls whether the Application attempts to
		// detect Docker.
		DetectDocker bool
		// DetectKubernetes controls whether the Application attempts to
		// detect Kubernetes.
		DetectKubernetes bool

		// These settings provide system information when custom values
		// are required.
		LogicalProcessors int
		TotalRAMMIB       int
		BillingHostname   string
	}

	// CrossApplicationTracer controls behaviour relating to cross application
	// tracing (CAT), available since Go Agent v0.11.  The CrossApplication
	// Tracer and the DistributedTracer cannot be simultaneously enabled.
	CrossApplicationTracer struct {
		Enabled bool
	}

	// DistributedTracer controls behaviour relating to Distributed Tracing,
	// available since Go Agent v2.1. The DistributedTracer and the
	// CrossApplicationTracer cannot be simultaneously enabled.
	DistributedTracer struct {
		Enabled bool
	}

	// SpanEvents controls behavior relating to Span Events.  Span Events
	// require that distributed tracing is enabled.
	SpanEvents struct {
		Enabled bool
	}

	// DatastoreTracer controls behavior relating to datastore segments.
	DatastoreTracer struct {
		InstanceReporting struct {
			Enabled bool
		}
		DatabaseNameReporting struct {
			Enabled bool
		}
		QueryParameters struct {
			Enabled bool
		}
		// SlowQuery controls the capture of slow query traces.  Slow
		// query traces show you instances of your slowest datastore
		// segments.
		SlowQuery struct {
			Enabled   bool
			Threshold time.Duration
		}
	}

	// Attributes controls the attributes included with errors and
	// transaction events.
	Attributes AttributeDestinationConfig

	// RuntimeSampler controls the collection of runtime statistics like
	// CPU/Memory usage, goroutine count, and GC pauses.
	RuntimeSampler struct {
		// Enabled controls whether runtime statistics are captured.
		Enabled bool
	}
}

// AttributeDestinationConfig controls the attributes included with errors and
// transaction events.
type AttributeDestinationConfig struct {
	Enabled bool
	Include []string
	Exclude []string
}

// NewConfig creates an Config populated with the given appname, license,
// and expected default values.
func NewConfig(appname, license string) Config {
	c := Config{}

	c.AppName = appname
	c.License = license
	c.Enabled = true
	c.Labels = make(map[string]string)
	c.CustomInsightsEvents.Enabled = true
	c.TransactionEvents.Enabled = true
	c.TransactionEvents.Attributes.Enabled = true
	c.HighSecurity = false
	c.ErrorCollector.Enabled = true
	c.ErrorCollector.CaptureEvents = true
	c.ErrorCollector.IgnoreStatusCodes = []int{
		http.StatusNotFound, // 404
	}
	c.ErrorCollector.Attributes.Enabled = true
	c.Utilization.DetectAWS = true
	c.Utilization.DetectAzure = true
	c.Utilization.DetectPCF = true
	c.Utilization.DetectGCP = true
	c.Utilization.DetectDocker = true
	c.Utilization.DetectKubernetes = true
	c.Attributes.Enabled = true
	c.RuntimeSampler.Enabled = true

	c.TransactionTracer.Enabled = true
	c.TransactionTracer.Threshold.IsApdexFailing = true
	c.TransactionTracer.Threshold.Duration = 500 * time.Millisecond
	c.TransactionTracer.SegmentThreshold = 2 * time.Millisecond
	c.TransactionTracer.StackTraceThreshold = 500 * time.Millisecond
	c.TransactionTracer.Attributes.Enabled = true

	c.CrossApplicationTracer.Enabled = true
	c.DistributedTracer.Enabled = false
	c.SpanEvents.Enabled = true

	c.DatastoreTracer.InstanceReporting.Enabled = true
	c.DatastoreTracer.DatabaseNameReporting.Enabled = true
	c.DatastoreTracer.QueryParameters.Enabled = true
	c.DatastoreTracer.SlowQuery.Enabled = true
	c.DatastoreTracer.SlowQuery.Threshold = 10 * time.Millisecond

	return c
}

const (
	licenseLength = 40
	appNameLimit  = 3
)

// The following errors will be returned if your Config fails to validate.
var (
	errLicenseLen                       = fmt.Errorf("license length is not %d", licenseLength)
	errAppNameMissing                   = errors.New("string AppName required")
	errAppNameLimit                     = fmt.Errorf("max of %d rollup application names", appNameLimit)
	errHighSecurityWithSecurityPolicies = errors.New("SecurityPoliciesToken and HighSecurity are incompatible; please ensure HighSecurity is set to false if SecurityPoliciesToken is a non-empty string and a security policy has been set for your account")
	errMixedTracers                     = errors.New("CrossApplicationTracer and DistributedTracer cannot be enabled simultaneously; please choose CrossApplicationTracer (available since v1.11) or DistributedTracer (available since v2.1)")
)

// Validate checks the config for improper fields.  If the config is invalid,
// newrelic.NewApplication returns an error.
func (c Config) Validate() error {
	if c.Enabled {
		if len(c.License) != licenseLength {
			return errLicenseLen
		}
	} else {
		// The License may be empty when the agent is not enabled.
		if len(c.License) != licenseLength && len(c.License) != 0 {
			return errLicenseLen
		}
	}
	if "" == c.AppName && c.Enabled {
		return errAppNameMissing
	}
	if c.HighSecurity && "" != c.SecurityPoliciesToken {
		return errHighSecurityWithSecurityPolicies
	}
	if c.CrossApplicationTracer.Enabled && c.DistributedTracer.Enabled {
		return errMixedTracers
	}
	if strings.Count(c.AppName, ";") >= appNameLimit {
		return errAppNameLimit
	}
	return nil
}
