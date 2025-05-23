package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoAwsProjectStackProps struct {
	awscdk.StackProps
}

func NewGoAwsProjectStack(scope constructs.Construct, id string, props *GoAwsProjectStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// create DB table here
	tableProps := &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("username"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName: jsii.String("userTable"),
	}
	table := awsdynamodb.NewTable(stack, jsii.String("myUserTable"), tableProps)

	// create lambda function
	lambdaProps := &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("main"),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("lambda/function.zip"), nil),
	}
	myFunction := awslambda.NewFunction(stack, jsii.String("myLambdaFunction"), lambdaProps)

	// grant lambda permissions to the table
	table.GrantReadWriteData(myFunction)

	// api gateway definition
	apiProps := &awsapigateway.RestApiProps{
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
			AllowMethods: jsii.Strings("GET", "POST", "DELETE", "PUT", "OPTIONS"),
			AllowOrigins: jsii.Strings("*"),
		},
		DeployOptions: &awsapigateway.StageOptions{
			LoggingLevel: awsapigateway.MethodLoggingLevel_INFO,
		},
		CloudWatchRole: jsii.Bool(true),
	}
	api := awsapigateway.NewRestApi(stack, jsii.String("myAPIGateway"), apiProps)

	// integrate api with lambda
	integration := awsapigateway.NewLambdaIntegration(myFunction, nil)

	// Routes
	registerRoute := api.Root().AddResource(jsii.String("register"), nil)
	registerRoute.AddMethod(jsii.String("POST"), integration, nil)
	loginRoute := api.Root().AddResource(jsii.String("login"), nil)
	loginRoute.AddMethod(jsii.String("POST"), integration, nil)
	protectedRoute := api.Root().AddResource(jsii.String("protected"), nil)
	protectedRoute.AddMethod(jsii.String("GET"), integration, nil)

	// jwt secret
	secretProps := &awssecretsmanager.SecretProps{
		SecretName:  jsii.String("jwtSigningSecret"),
		Description: jsii.String("Secret key for signing JWT tokens"),
		GenerateSecretString: &awssecretsmanager.SecretStringGenerator{
			SecretStringTemplate: jsii.String(`{"key":""}`),
			GenerateStringKey:    jsii.String("key"),
			PasswordLength:       jsii.Number(32),
		},
	}
	jwtSecret := awssecretsmanager.NewSecret(stack, jsii.String("myJWTSecret"), secretProps)

	// grant secret read access to lambda
	jwtSecret.GrantRead(myFunction.Role(), &[]*string{jsii.String("AWSCURRENT")})

	// pass the secret ARN as env to lambda
	myFunction.AddEnvironment(jsii.String("JWT_SECRET_ARN"), jwtSecret.SecretArn(), nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoAwsProjectStack(app, "GoAwsProjectStack", &GoAwsProjectStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
