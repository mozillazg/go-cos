package cos

import (
	"testing"
)

func TestBucketService_GetACL(t *testing.T) {

	setup()
	defer teardown()

	//	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//		testMethod(t, r, "GET")
	//		vs := values{
	//			"acl": "",
	//		}
	//		testFormValues(t, r, vs)
	//		fmt.Fprint(w, `<AccessControlPolicy>
	//	<Owner>
	//		<uin>100000760461</uin>
	//	</Owner>
	//	<AccessControlList>
	//		<Grant>
	//			<Grantee type="RootAccount">
	//				<uin>100000760461</uin>
	//			</Grantee>
	//			<Permission>FULL_CONTROL</Permission>
	//		</Grant>
	//		<Grant>
	//			<Grantee type="RootAccount">
	//				<uin>100000760461</uin>
	//			</Grantee>
	//			<Permission>FULL_CONTROL</Permission>
	//		</Grant>
	//	</AccessControlList>
	//</AccessControlPolicy>`)
	//	})
	//
	//	ref, _, err := client.Bucket.GetACL(context.Background(), NewAuthTime(time.Minute))
	//	if err != nil {
	//		t.Fatalf("Bucket.GetACL returned error: %v", err)
	//	}
	//
	//	want := &BucketGetACLResult{
	//		XMLName:  xml.Name{Local: "Tagging"},
	//		TagSet: []BucketTaggingTag{
	//			{"test_k2", "test_v2"},
	//			{"test_k3", "test_vv"},
	//		},
	//	}
	//
	//	if !reflect.DeepEqual(ref, want) {
	//		t.Errorf("Bucket.GetACL returned %+v, want %+v", ref, want)
	//	}

}

func TestBucketService_PutACL(t *testing.T) {

}
