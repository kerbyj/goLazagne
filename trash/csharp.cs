using System;
using System.Runtime.InteropServices;
public class getData
{
    [StructLayout(LayoutKind.Sequential, CharSet = CharSet.Unicode)]
    public struct NativeCredential
    {
        public uint Flags;
        public enum CredType : uint{
            Generic = 1,
            DomainPassword = 2,
            DomainCertificate = 3,
            DomainVisiblePassword = 4,
            GenericCertificate = 5,
            DomainExtended = 6,
            Maximum = 7,
            MaximumEx = (Maximum + 1000)
        }
        public IntPtr TargetName;
        public IntPtr Comment;
        public System.Runtime.InteropServices.ComTypes.FILETIME LastWritten;
        public uint CredentialBlobSize;
        public IntPtr CredentialBlob;
        public uint Persist;
        public uint AttributeCount;
        public IntPtr Attributes;
        public IntPtr TargetAlias;
        public IntPtr UserName;
    }
    [DllImport("Advapi32.dll", SetLastError = true, EntryPoint = "CredEnumerateA", CharSet = CharSet.Unicode)]
    public static extern bool CredEnumerate([In] string filter, [In] int flags, out int count, out IntPtr credentialPtrs);
    public static void get(){
        int count;
        IntPtr pCredentials;
        bool read = CredEnumerate(null, 0x0, out count, out pCredentials);
        for (int inx = 0; inx < count; inx++)
        {
            IntPtr pCred = Marshal.ReadIntPtr(pCredentials, inx * IntPtr.Size);
            NativeCredential nativeCredential = (NativeCredential)Marshal.PtrToStructure(pCred, typeof(NativeCredential)); // Native
            string username = Marshal.PtrToStringUni(nativeCredential.TargetName);
            if (0 < nativeCredential.CredentialBlobSize){
                string targetname = Marshal.PtrToStringAnsi(nativeCredential.TargetName);
                string password = Marshal.PtrToStringUni(nativeCredential.CredentialBlob, (int)nativeCredential.CredentialBlobSize / 2);
                if(password.Length > 64) {continue;}
                Console.WriteLine(targetname + " " + password);
            }
        }
    }
}